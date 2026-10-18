package services

import (
	"bytes"
	"context"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"os"
	"strings"
	"testing"
	"time"

	"opendomain/internal/config"
	"opendomain/pkg/logger"
)

func TestMain(m *testing.M) {
	logger.Init("error")
	os.Exit(m.Run())
}

func TestBuildMessageEncodesSubjectAndBothParts(t *testing.T) {
	subject := "[Yunamo] 验证你的邮箱地址 - 这是一个比较长的中文主题，用来测试折行是否正确处理"
	raw, err := BuildMessage("Yunamo", "noreply@yunamo.com", &Mail{
		To:      "user@example.com",
		Subject: subject,
		Text:    "纯文本 body",
		HTML:    "<p>HTML 正文</p>",
	}, time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC))
	if err != nil {
		t.Fatalf("BuildMessage: %v", err)
	}

	msg, err := mail.ReadMessage(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("ReadMessage: %v", err)
	}
	if got := msg.Header.Get("Subject"); !strings.HasPrefix(got, "=?UTF-8?b?") {
		t.Fatalf("subject not RFC 2047 encoded: %q", got)
	}
	decoded, err := new(mime.WordDecoder).DecodeHeader(msg.Header.Get("Subject"))
	if err != nil || decoded != subject {
		t.Fatalf("decoded subject = %q, %v; want %q", decoded, err, subject)
	}
	if got := msg.Header.Get("From"); got != `"Yunamo" <noreply@yunamo.com>` {
		t.Fatalf("From = %q", got)
	}
	if got := msg.Header.Get("To"); got != "<user@example.com>" {
		t.Fatalf("To = %q", got)
	}
	if got := msg.Header.Get("Message-Id"); !strings.HasPrefix(got, "<") || !strings.HasSuffix(got, "@yunamo.com>") {
		t.Fatalf("Message-ID = %q", got)
	}
	if msg.Header.Get("Date") == "" || msg.Header.Get("Mime-Version") != "1.0" {
		t.Fatalf("missing Date or MIME-Version: %v", msg.Header)
	}

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil || mediaType != "multipart/alternative" {
		t.Fatalf("Content-Type = %q, %v", msg.Header.Get("Content-Type"), err)
	}
	mr := multipart.NewReader(msg.Body, params["boundary"])
	want := []struct{ contentType, body string }{
		{"text/plain; charset=UTF-8", "纯文本 body"},
		{"text/html; charset=UTF-8", "<p>HTML 正文</p>"},
	}
	for i, w := range want {
		p, err := mr.NextPart()
		if err != nil {
			t.Fatalf("part %d: %v", i, err)
		}
		if got := p.Header.Get("Content-Type"); got != w.contentType {
			t.Fatalf("part %d Content-Type = %q", i, got)
		}
		body, _ := io.ReadAll(p)
		if string(body) != w.body {
			t.Fatalf("part %d body = %q, want %q", i, body, w.body)
		}
	}
	if _, err := mr.NextPart(); err != io.EOF {
		t.Fatalf("expected exactly two parts, got err %v", err)
	}
}

func TestBuildMessageRejectsHeaderInjection(t *testing.T) {
	cases := []Mail{
		{To: "user@example.com", Subject: "hi\r\nBcc: victim@example.com"},
		{To: "user@example.com", Subject: "hi\nBcc: victim@example.com"},
		{To: "user@example.com\r\nBcc: victim@example.com", Subject: "hi"},
		{To: "user@example.com, victim@example.com", Subject: "hi"},
		{To: "Name <user@example.com>", Subject: "hi"},
		{To: "user@example.com", ReplyTo: "a@b.co\r\nBcc: victim@example.com", Subject: "hi"},
		{To: "user@example.com", ReplyTo: "a@b.co, victim@example.com", Subject: "hi"},
	}
	for _, m := range cases {
		m := m
		if _, err := BuildMessage("Yunamo", "noreply@yunamo.com", &m, time.Now()); err == nil {
			t.Errorf("BuildMessage(%q, %q) succeeded, want error", m.To, m.Subject)
		}
	}
	if _, err := BuildMessage("Yunamo\r\nX: y", "noreply@yunamo.com", &Mail{To: "a@b.co", Subject: "hi"}, time.Now()); err == nil {
		t.Error("expected error for CR/LF in sender name")
	}
}

func TestBuildMessageReplyTo(t *testing.T) {
	raw, err := BuildMessage("Yunamo", "noreply@yunamo.com", &Mail{To: "support@yunamo.com", ReplyTo: "ada@example.com", Subject: "hi"}, time.Now())
	if err != nil {
		t.Fatalf("BuildMessage: %v", err)
	}
	msg, err := mail.ReadMessage(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("ReadMessage: %v", err)
	}
	if got := msg.Header.Get("Reply-To"); got != "<ada@example.com>" {
		t.Fatalf("Reply-To = %q", got)
	}
	raw, _ = BuildMessage("Yunamo", "noreply@yunamo.com", &Mail{To: "a@b.co", Subject: "hi"}, time.Now())
	if strings.Contains(string(raw), "Reply-To:") {
		t.Fatal("Reply-To header written without a reply address")
	}
}

func TestSMTPMailerSkipsWithoutHost(t *testing.T) {
	s := NewSMTPMailer(&config.Config{SiteName: "Yunamo"})
	if err := s.Send(context.Background(), &Mail{To: "a@b.co", Subject: "x"}); err != nil {
		t.Fatalf("Send without MAIL_HOST = %v, want nil", err)
	}
}

func TestNormalizeLocale(t *testing.T) {
	for in, want := range map[string]string{"": LocaleZhCN, "zh-CN": LocaleZhCN, "en-US": LocaleEnUS, "en": LocaleEnUS, "fr": LocaleZhCN} {
		if got := NormalizeLocale(in); got != want {
			t.Errorf("NormalizeLocale(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestTemplatesEscapeUserValuesAndFollowLocale(t *testing.T) {
	link := "https://yunamo.com/verify-email?token=abc"
	zh, err := VerifyEmailMail("zh-CN", "Yunamo", "a@b.co", `<script>alert(1)</script>`, link)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(zh.HTML, "<script>") {
		t.Fatal("username was not escaped in HTML")
	}
	if !strings.Contains(zh.Subject, "验证") || !strings.Contains(zh.HTML, "你的动力源") || !strings.Contains(zh.Text, link) {
		t.Fatalf("unexpected zh-CN mail: %q / %q", zh.Subject, zh.Text)
	}
	if !strings.Contains(zh.HTML, "#3EE6B4") || !strings.Contains(zh.HTML, "#04110C") {
		t.Fatal("button colours missing")
	}

	en, err := ResetPasswordMail("en-US", "Yunamo", "a@b.co", "bob", "https://yunamo.com/reset-password?token=abc")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(en.Subject, "Reset your password") || !strings.Contains(en.Text, "Your dynamo") {
		t.Fatalf("unexpected en-US mail: %q / %q", en.Subject, en.Text)
	}

	alert, err := AbuseAlertMail("Yunamo", "admin@b.co", AbuseAlert{
		ID: 7, Category: "phishing", Subdomain: "foo.loc.cc", Details: strings.Repeat("x", 3000) + "<b>",
		CreatedAt: time.Now(), AdminURL: "https://yunamo.com/admin/abuse-reports",
	})
	if err != nil {
		t.Fatal(err)
	}
	if alert.Subject != "[Yunamo] New abuse report #7: phishing · foo.loc.cc" {
		t.Fatalf("alert subject = %q", alert.Subject)
	}
	if strings.Contains(alert.Text, strings.Repeat("x", 1001)) {
		t.Fatal("details were not truncated")
	}
}
