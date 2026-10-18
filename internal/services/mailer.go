package services

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strconv"
	"strings"
	"time"

	"opendomain/internal/config"
	"opendomain/pkg/logger"
)

const smtpTimeout = 20 * time.Second

// Mail is a single outgoing message with a plain-text and an HTML body.
type Mail struct {
	To      string
	ReplyTo string
	Subject string
	Text    string
	HTML    string
}

// MailSender delivers a message. Handlers depend on this interface so tests
// can substitute a fake transport.
type MailSender interface {
	Send(ctx context.Context, m *Mail) error
}

// SMTPMailer sends mail through the server configured by MAIL_* settings.
type SMTPMailer struct {
	cfg *config.Config
}

func NewSMTPMailer(cfg *config.Config) *SMTPMailer {
	return &SMTPMailer{cfg: cfg}
}

var errUnsafeHeader = errors.New("mail: header value contains CR or LF")

func hasLineBreak(s string) bool {
	return strings.ContainsAny(s, "\r\n")
}

// parseBareAddress accepts only a plain addr-spec such as a@b.co, so a value
// like "x@y.z, other@victim" cannot add recipients.
func parseBareAddress(s string) (string, error) {
	if hasLineBreak(s) {
		return "", errUnsafeHeader
	}
	addr, err := mail.ParseAddress(s)
	if err != nil {
		return "", fmt.Errorf("mail: invalid address %q: %w", s, err)
	}
	if addr.Address != s {
		return "", fmt.Errorf("mail: invalid address %q", s)
	}
	return addr.Address, nil
}

func addressDomain(addr string) string {
	if at := strings.LastIndex(addr, "@"); at >= 0 && at+1 < len(addr) {
		return addr[at+1:]
	}
	return "localhost"
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// encodeSubject RFC 2047-encodes non-ASCII subjects and folds between
// encoded words so long Chinese subjects stay within line limits.
func encodeSubject(subject string) string {
	enc := mime.BEncoding.Encode("UTF-8", subject)
	return strings.ReplaceAll(enc, "?= =?", "?=\r\n =?")
}

func writeQuotedPrintablePart(w *multipart.Writer, contentType, body string) error {
	h := textproto.MIMEHeader{}
	h.Set("Content-Type", contentType)
	h.Set("Content-Transfer-Encoding", "quoted-printable")
	pw, err := w.CreatePart(h)
	if err != nil {
		return err
	}
	qp := quotedprintable.NewWriter(pw)
	if _, err := qp.Write([]byte(body)); err != nil {
		return err
	}
	return qp.Close()
}

// BuildMessage renders a complete RFC 5322 message with a
// multipart/alternative body.
func BuildMessage(fromName, fromAddr string, m *Mail, now time.Time) ([]byte, error) {
	if hasLineBreak(fromName) || hasLineBreak(m.Subject) {
		return nil, errUnsafeHeader
	}
	from, err := parseBareAddress(fromAddr)
	if err != nil {
		return nil, err
	}
	to, err := parseBareAddress(m.To)
	if err != nil {
		return nil, err
	}
	replyTo := ""
	if m.ReplyTo != "" {
		if replyTo, err = parseBareAddress(m.ReplyTo); err != nil {
			return nil, err
		}
	}
	msgID, err := randomHex(16)
	if err != nil {
		return nil, err
	}

	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	if err := writeQuotedPrintablePart(mw, "text/plain; charset=UTF-8", m.Text); err != nil {
		return nil, err
	}
	if err := writeQuotedPrintablePart(mw, "text/html; charset=UTF-8", m.HTML); err != nil {
		return nil, err
	}
	if err := mw.Close(); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	header := func(k, v string) {
		buf.WriteString(k)
		buf.WriteString(": ")
		buf.WriteString(v)
		buf.WriteString("\r\n")
	}
	header("From", (&mail.Address{Name: fromName, Address: from}).String())
	header("To", (&mail.Address{Address: to}).String())
	if replyTo != "" {
		header("Reply-To", (&mail.Address{Address: replyTo}).String())
	}
	header("Subject", encodeSubject(m.Subject))
	header("Date", now.Format(time.RFC1123Z))
	header("Message-ID", "<"+msgID+"@"+addressDomain(from)+">")
	header("MIME-Version", "1.0")
	header("Content-Type", `multipart/alternative; boundary="`+mw.Boundary()+`"`)
	buf.WriteString("\r\n")
	buf.Write(body.Bytes())
	return buf.Bytes(), nil
}

// Send delivers m over SMTP. Port 465 uses implicit TLS; any other port
// requires STARTTLS so credentials never cross the network in clear text.
func (s *SMTPMailer) Send(ctx context.Context, m *Mail) error {
	ec := s.cfg.Email
	if ec.Host == "" {
		logger.Warnf("mail: MAIL_HOST is not configured, skipping message %q to %s", m.Subject, m.To)
		return nil
	}
	fromAddr := ec.From
	if fromAddr == "" {
		fromAddr = ec.User
	}
	msg, err := BuildMessage(s.cfg.SiteName, fromAddr, m, time.Now())
	if err != nil {
		return err
	}

	port := ec.Port
	if port == 0 {
		port = 587
	}
	addr := net.JoinHostPort(ec.Host, strconv.Itoa(port))
	tlsCfg := &tls.Config{ServerName: ec.Host, MinVersion: tls.VersionTLS12}
	dialer := &net.Dialer{Timeout: smtpTimeout}

	ctx, cancel := context.WithTimeout(ctx, smtpTimeout)
	defer cancel()

	var conn net.Conn
	if port == 465 {
		conn, err = (&tls.Dialer{NetDialer: dialer, Config: tlsCfg}).DialContext(ctx, "tcp", addr)
	} else {
		conn, err = dialer.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return fmt.Errorf("mail: dial %s: %w", addr, err)
	}
	deadline, _ := ctx.Deadline()
	_ = conn.SetDeadline(deadline)

	c, err := smtp.NewClient(conn, ec.Host)
	if err != nil {
		conn.Close()
		return fmt.Errorf("mail: smtp handshake: %w", err)
	}
	defer c.Close()

	if err := c.Hello(addressDomain(fromAddr)); err != nil {
		return fmt.Errorf("mail: EHLO: %w", err)
	}
	if port != 465 {
		if ok, _ := c.Extension("STARTTLS"); !ok {
			return errors.New("mail: server does not offer STARTTLS, refusing to send")
		}
		if err := c.StartTLS(tlsCfg); err != nil {
			return fmt.Errorf("mail: STARTTLS: %w", err)
		}
	}
	if ec.User != "" {
		if err := c.Auth(smtp.PlainAuth("", ec.User, ec.Password, ec.Host)); err != nil {
			return fmt.Errorf("mail: AUTH: %w", err)
		}
	}
	if err := c.Mail(fromAddr); err != nil {
		return fmt.Errorf("mail: MAIL FROM: %w", err)
	}
	if err := c.Rcpt(m.To); err != nil {
		return fmt.Errorf("mail: RCPT TO: %w", err)
	}
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("mail: DATA: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		w.Close()
		return fmt.Errorf("mail: write body: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("mail: end DATA: %w", err)
	}
	return c.Quit()
}

// RunAsync runs fn in a goroutine and logs instead of crashing on panic.
func RunAsync(name string, fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("%s: panic: %v", name, r)
			}
		}()
		fn()
	}()
}

// Mailer queues mail and other mail-related work in the background so
// requests never wait on SMTP.
type Mailer struct {
	Sender MailSender
	// Run executes background work; tests swap in a synchronous runner.
	Run func(name string, fn func())
}

func NewMailer(sender MailSender) *Mailer {
	return &Mailer{Sender: sender, Run: RunAsync}
}

// Go runs fn in the background.
func (m *Mailer) Go(name string, fn func()) {
	if m == nil {
		return
	}
	run := m.Run
	if run == nil {
		run = RunAsync
	}
	run(name, fn)
}

// SendAsync delivers msg in the background and logs failures.
func (m *Mailer) SendAsync(msg *Mail) {
	if m == nil || m.Sender == nil || msg == nil {
		return
	}
	m.Go("mail", func() {
		if err := m.Sender.Send(context.Background(), msg); err != nil {
			logger.Errorf("mail: failed to send %q to %s: %v", msg.Subject, msg.To, err)
		}
	})
}
