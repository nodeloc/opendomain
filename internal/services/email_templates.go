package services

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"strings"
	texttemplate "text/template"
	"time"
	"unicode/utf8"
)

const (
	LocaleZhCN = "zh-CN"
	LocaleEnUS = "en-US"

	abuseAlertDetailsMax = 1000
)

// NormalizeLocale maps any locale string onto the two supported email
// languages, defaulting to Chinese like the frontend does.
func NormalizeLocale(locale string) string {
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(locale)), "en") {
		return LocaleEnUS
	}
	return LocaleZhCN
}

type emailRow struct {
	Label string
	Value string
}

type emailLayout struct {
	Lang       string
	SiteName   string
	Preheader  string
	Heading    string
	Subheading string
	Greeting   string
	Paragraphs []string
	Rows       []emailRow
	QuoteLabel string
	Quote      string
	ButtonText string
	ButtonURL  string
	LinkHint   string
	Notes      []string
	Footer     string
}

const emailHTMLSource = `<!DOCTYPE html>
<html lang="{{.Lang}}">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Heading}}</title>
</head>
<body style="margin:0;padding:0;background-color:#F4F6F5;color:#0F1A15;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,'PingFang SC','Hiragino Sans GB','Microsoft YaHei','Helvetica Neue',Arial,sans-serif;">
<div style="display:none;max-height:0;overflow:hidden;opacity:0;">{{.Preheader}}</div>
<table role="presentation" width="100%" cellspacing="0" cellpadding="0" border="0" style="background-color:#F4F6F5;">
<tr><td align="center" style="padding:32px 16px;">
<table role="presentation" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width:560px;">
<tr><td style="padding:0 4px 20px 4px;font-size:22px;font-weight:700;letter-spacing:-0.02em;color:#0F1A15;">{{.SiteName}}</td></tr>
<tr><td style="background-color:#FFFFFF;border:1px solid #E2E8E5;border-radius:12px;padding:32px 28px;">
<h1 style="margin:0 0 8px 0;font-size:20px;line-height:1.4;font-weight:700;color:#0F1A15;">{{.Heading}}</h1>
{{- if .Subheading}}
<p style="margin:0 0 16px 0;font-size:14px;line-height:1.5;color:#5B6B64;">{{.Subheading}}</p>
{{- end}}
{{- if .Greeting}}
<p style="margin:16px 0 0 0;font-size:15px;line-height:1.6;color:#0F1A15;">{{.Greeting}}</p>
{{- end}}
{{- range .Paragraphs}}
<p style="margin:12px 0 0 0;font-size:15px;line-height:1.6;color:#0F1A15;">{{.}}</p>
{{- end}}
{{- if .Rows}}
<table role="presentation" width="100%" cellspacing="0" cellpadding="0" border="0" style="margin:20px 0 0 0;border-collapse:collapse;">
{{- range .Rows}}
<tr>
<td style="padding:8px 12px 8px 0;border-bottom:1px solid #EEF2F0;font-size:13px;line-height:1.5;color:#5B6B64;white-space:nowrap;vertical-align:top;">{{.Label}}</td>
<td style="padding:8px 0;border-bottom:1px solid #EEF2F0;font-size:14px;line-height:1.5;color:#0F1A15;word-break:break-all;vertical-align:top;">{{.Value}}</td>
</tr>
{{- end}}
</table>
{{- end}}
{{- if .Quote}}
<p style="margin:20px 0 6px 0;font-size:13px;line-height:1.5;color:#5B6B64;">{{.QuoteLabel}}</p>
<div style="margin:0;padding:12px 14px;background-color:#F4F6F5;border-radius:8px;font-size:14px;line-height:1.6;color:#0F1A15;white-space:pre-wrap;word-break:break-word;">{{.Quote}}</div>
{{- end}}
{{- if .ButtonURL}}
<table role="presentation" cellspacing="0" cellpadding="0" border="0" style="margin:28px 0 0 0;">
<tr><td style="border-radius:8px;background-color:#3EE6B4;">
<a href="{{.ButtonURL}}" target="_blank" style="display:inline-block;padding:12px 24px;font-size:15px;font-weight:600;color:#04110C;text-decoration:none;border-radius:8px;">{{.ButtonText}}</a>
</td></tr>
</table>
{{- if .LinkHint}}
<p style="margin:20px 0 4px 0;font-size:13px;line-height:1.5;color:#5B6B64;">{{.LinkHint}}</p>
<p style="margin:0;font-size:13px;line-height:1.5;word-break:break-all;"><a href="{{.ButtonURL}}" style="color:#0F7A5C;">{{.ButtonURL}}</a></p>
{{- end}}
{{- end}}
{{- range .Notes}}
<p style="margin:16px 0 0 0;font-size:13px;line-height:1.5;color:#5B6B64;">{{.}}</p>
{{- end}}
</td></tr>
<tr><td style="padding:20px 4px 0 4px;font-size:12px;line-height:1.5;color:#7A8A83;">{{.Footer}}</td></tr>
</table>
</td></tr>
</table>
</body>
</html>
`

const emailTextSource = `{{.Heading}}
{{- if .Subheading}}
{{.Subheading}}
{{- end}}
{{if .Greeting}}
{{.Greeting}}
{{end}}
{{- range .Paragraphs}}
{{.}}
{{end}}
{{- if .Rows}}
{{range .Rows}}{{.Label}}: {{.Value}}
{{end}}
{{- end}}
{{- if .Quote}}
{{.QuoteLabel}}
{{.Quote}}
{{end}}
{{- if .ButtonURL}}
{{.ButtonText}}: {{.ButtonURL}}
{{end}}
{{- range .Notes}}
{{.}}
{{end}}
--
{{.Footer}}
`

var (
	emailHTMLTemplate = htmltemplate.Must(htmltemplate.New("email.html").Parse(emailHTMLSource))
	emailTextTemplate = texttemplate.Must(texttemplate.New("email.txt").Parse(emailTextSource))
)

func siteNameOrDefault(siteName string) string {
	if s := strings.TrimSpace(siteName); s != "" {
		return s
	}
	return "Yunamo"
}

func footerFor(locale, siteName string) string {
	if locale == LocaleEnUS {
		return siteName + " · Your dynamo"
	}
	return siteName + " · 你的动力源"
}

func renderMail(to, subject string, layout *emailLayout) (*Mail, error) {
	var html, text bytes.Buffer
	if err := emailHTMLTemplate.Execute(&html, layout); err != nil {
		return nil, fmt.Errorf("render html email: %w", err)
	}
	if err := emailTextTemplate.Execute(&text, layout); err != nil {
		return nil, fmt.Errorf("render text email: %w", err)
	}
	return &Mail{To: to, Subject: subject, Text: text.String(), HTML: html.String()}, nil
}

func truncateRunes(s string, max int) string {
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	r := []rune(s)
	return string(r[:max]) + "…"
}

// singleLine keeps user-supplied values that end up in a subject line free of
// line breaks.
func singleLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// VerifyEmailMail builds the "confirm your email address" message.
func VerifyEmailMail(locale, siteName, to, username, link string) (*Mail, error) {
	locale = NormalizeLocale(locale)
	siteName = siteNameOrDefault(siteName)
	l := &emailLayout{Lang: locale, SiteName: siteName, ButtonURL: link, Footer: footerFor(locale, siteName)}
	var subject string
	if locale == LocaleEnUS {
		subject = fmt.Sprintf("[%s] Verify your email address", siteName)
		l.Heading = "Verify your email address"
		l.Preheader = "Confirm your email to start claiming domains."
		if username != "" {
			l.Greeting = "Hi " + username + ","
		}
		l.Paragraphs = []string{
			"Thanks for signing up for " + siteName + ". Confirm that this is your email address to start claiming domains.",
		}
		l.ButtonText = "Verify email"
		l.LinkHint = "If the button doesn't work, copy this link into your browser:"
		l.Notes = []string{
			"This link expires in 24 hours and can only be used once.",
			"If you didn't create an account, you can ignore this email.",
		}
	} else {
		subject = fmt.Sprintf("[%s] 验证你的邮箱地址", siteName)
		l.Heading = "验证你的邮箱地址"
		l.Preheader = "确认邮箱后即可开始领取域名。"
		if username != "" {
			l.Greeting = username + "，你好："
		}
		l.Paragraphs = []string{
			"感谢注册 " + siteName + "。请确认这是你的邮箱地址，验证后即可开始领取域名。",
		}
		l.ButtonText = "验证邮箱"
		l.LinkHint = "如果按钮无法点击，请将下面的链接复制到浏览器中打开："
		l.Notes = []string{
			"该链接 24 小时内有效，且只能使用一次。",
			"如果这不是你本人的操作，请忽略本邮件。",
		}
	}
	return renderMail(to, subject, l)
}

// ResetPasswordMail builds the password reset message.
func ResetPasswordMail(locale, siteName, to, username, link string) (*Mail, error) {
	locale = NormalizeLocale(locale)
	siteName = siteNameOrDefault(siteName)
	l := &emailLayout{Lang: locale, SiteName: siteName, ButtonURL: link, Footer: footerFor(locale, siteName)}
	var subject string
	if locale == LocaleEnUS {
		subject = fmt.Sprintf("[%s] Reset your password", siteName)
		l.Heading = "Reset your password"
		l.Preheader = "Use this link to choose a new password."
		if username != "" {
			l.Greeting = "Hi " + username + ","
		}
		l.Paragraphs = []string{
			"We received a request to reset the password for your " + siteName + " account. Click the button below to choose a new one.",
		}
		l.ButtonText = "Reset password"
		l.LinkHint = "If the button doesn't work, copy this link into your browser:"
		l.Notes = []string{
			"This link expires in 1 hour and can only be used once.",
			"If you didn't ask to reset your password, you can ignore this email. Your password won't change.",
		}
	} else {
		subject = fmt.Sprintf("[%s] 重置你的密码", siteName)
		l.Heading = "重置你的密码"
		l.Preheader = "使用此链接设置新密码。"
		if username != "" {
			l.Greeting = username + "，你好："
		}
		l.Paragraphs = []string{
			"我们收到了重置你 " + siteName + " 账户密码的请求。点击下方按钮设置新密码。",
		}
		l.ButtonText = "重置密码"
		l.LinkHint = "如果按钮无法点击，请将下面的链接复制到浏览器中打开："
		l.Notes = []string{
			"该链接 1 小时内有效，且只能使用一次。",
			"如果这不是你本人的操作，请忽略本邮件，你的密码不会被更改。",
		}
	}
	return renderMail(to, subject, l)
}

// AbuseAlert carries the fields shown to admins for a new abuse report.
type AbuseAlert struct {
	ID            uint
	Category      string
	Subdomain     string
	URL           string
	DomainStatus  string
	Details       string
	ReporterEmail string
	CreatedAt     time.Time
	AdminURL      string
}

// AbuseAlertMail builds the admin notification for a new abuse report. It is
// bilingual because admins do not share a single locale.
func AbuseAlertMail(siteName, to string, a AbuseAlert) (*Mail, error) {
	siteName = siteNameOrDefault(siteName)
	subject := fmt.Sprintf("[%s] New abuse report #%d: %s · %s", siteName, a.ID, singleLine(a.Category), singleLine(a.Subdomain))
	or := func(v, empty string) string {
		if strings.TrimSpace(v) == "" {
			return empty
		}
		return v
	}
	l := &emailLayout{
		Lang:       LocaleZhCN,
		SiteName:   siteName,
		Preheader:  fmt.Sprintf("#%d %s · %s", a.ID, a.Category, a.Subdomain),
		Heading:    fmt.Sprintf("收到新的滥用举报 #%d", a.ID),
		Subheading: fmt.Sprintf("New abuse report #%d", a.ID),
		Paragraphs: []string{"有人提交了一条新的滥用举报，请尽快处理。 / A new abuse report was submitted and needs review."},
		Rows: []emailRow{
			{Label: "类别 Category", Value: a.Category},
			{Label: "域名 Subdomain", Value: a.Subdomain},
			{Label: "URL", Value: or(a.URL, "-")},
			{Label: "域名状态 Domain status", Value: or(a.DomainStatus, "未注册 / not registered here")},
			{Label: "举报人邮箱 Reporter email", Value: or(a.ReporterEmail, "-")},
			{Label: "提交时间 Submitted", Value: a.CreatedAt.UTC().Format("2006-01-02 15:04:05 UTC")},
		},
		QuoteLabel: "详细说明 Details",
		Quote:      truncateRunes(a.Details, abuseAlertDetailsMax),
		ButtonText: "处理举报 Review reports",
		ButtonURL:  a.AdminURL,
		Footer:     siteName + " · 你的动力源 / Your dynamo",
	}
	return renderMail(to, subject, l)
}

// AbuseOutcomeMail tells a reporter how their report was closed. It must
// never include internal admin notes.
func AbuseOutcomeMail(locale, siteName, to string, reportID uint, subdomain string, resolved bool) (*Mail, error) {
	locale = NormalizeLocale(locale)
	siteName = siteNameOrDefault(siteName)
	l := &emailLayout{Lang: locale, SiteName: siteName, Footer: footerFor(locale, siteName)}
	var subject string
	if locale == LocaleEnUS {
		l.Rows = []emailRow{{Label: "Report", Value: fmt.Sprintf("#%d", reportID)}, {Label: "Domain", Value: subdomain}}
		l.Notes = []string{"You're receiving this because you left your email address on an abuse report."}
		if resolved {
			subject = fmt.Sprintf("[%s] We took action on your report #%d", siteName, reportID)
			l.Heading = "We took action on " + subdomain
			l.Paragraphs = []string{"Thanks for your report. We reviewed it and took action on " + subdomain + "."}
		} else {
			subject = fmt.Sprintf("[%s] Update on your report #%d", siteName, reportID)
			l.Heading = "We reviewed your report"
			l.Paragraphs = []string{"Thanks for your report about " + subdomain + ". We reviewed it and found no violation of our policies, so no action was taken."}
		}
		l.Preheader = l.Paragraphs[0]
	} else {
		l.Rows = []emailRow{{Label: "举报编号", Value: fmt.Sprintf("#%d", reportID)}, {Label: "域名", Value: subdomain}}
		l.Notes = []string{"你收到这封邮件，是因为你在提交滥用举报时留下了邮箱地址。"}
		if resolved {
			subject = fmt.Sprintf("[%s] 我们已处理你的举报 #%d", siteName, reportID)
			l.Heading = "我们已对 " + subdomain + " 采取措施"
			l.Paragraphs = []string{"感谢你的举报。我们已完成审核，并对 " + subdomain + " 采取了相应措施。"}
		} else {
			subject = fmt.Sprintf("[%s] 你的举报 #%d 处理结果", siteName, reportID)
			l.Heading = "我们已审核你的举报"
			l.Paragraphs = []string{"感谢你对 " + subdomain + " 的举报。我们审核后未发现违反我们政策的情况，因此没有采取措施。"}
		}
		l.Preheader = l.Paragraphs[0]
	}
	return renderMail(to, subject, l)
}

// PartnerApplicationAlert carries the fields of a new partnership application.
type PartnerApplicationAlert struct {
	ID              uint
	Name            string
	Email           string
	Organization    string
	Website         string
	PartnershipType string
	Contact         string
	Message         string
	Locale          string
	CreatedAt       time.Time
	AdminURL        string
}

var partnershipTypeLabels = map[string]string{
	"sponsorship": "资金赞助 Sponsorship",
	"resources":   "资源支持 Resources",
	"business":    "商务合作 Business",
	"community":   "社区与媒体 Community & media",
	"other":       "其他 Other",
}

// PartnerApplicationMail builds the notification sent to the support inbox
// for a new partnership application. Reply-To is the applicant, so replying
// from the inbox answers them directly.
func PartnerApplicationMail(siteName, to string, a PartnerApplicationAlert) (*Mail, error) {
	siteName = siteNameOrDefault(siteName)
	who := singleLine(a.Name)
	if org := singleLine(a.Organization); org != "" {
		who += " · " + org
	}
	typeLabel := partnershipTypeLabels[a.PartnershipType]
	if typeLabel == "" {
		typeLabel = a.PartnershipType
	}
	subject := fmt.Sprintf("[%s] New partner application #%d: %s", siteName, a.ID, truncateRunes(who, 80))
	or := func(v string) string {
		if strings.TrimSpace(v) == "" {
			return "-"
		}
		return v
	}
	l := &emailLayout{
		Lang:       LocaleZhCN,
		SiteName:   siteName,
		Preheader:  fmt.Sprintf("#%d %s · %s", a.ID, who, typeLabel),
		Heading:    fmt.Sprintf("收到新的合作申请 #%d", a.ID),
		Subheading: fmt.Sprintf("New partner application #%d", a.ID),
		Paragraphs: []string{"有人通过赞助页面提交了合作申请。直接回复本邮件即可联系对方。 / Someone applied to partner with us from the sponsor page. Reply to this email to answer them directly."},
		Rows: []emailRow{
			{Label: "合作类型 Type", Value: typeLabel},
			{Label: "联系人 Name", Value: a.Name},
			{Label: "邮箱 Email", Value: a.Email},
			{Label: "组织 Organization", Value: or(a.Organization)},
			{Label: "网站 Website", Value: or(a.Website)},
			{Label: "其他联系方式 Contact", Value: or(a.Contact)},
			{Label: "语言 Locale", Value: or(a.Locale)},
			{Label: "提交时间 Submitted", Value: a.CreatedAt.UTC().Format("2006-01-02 15:04:05 UTC")},
		},
		QuoteLabel: "合作说明 Message",
		Quote:      a.Message,
		ButtonText: "查看合作申请 View applications",
		ButtonURL:  a.AdminURL,
		Footer:     siteName + " · 你的动力源 / Your dynamo",
	}
	m, err := renderMail(to, subject, l)
	if err != nil {
		return nil, err
	}
	m.ReplyTo = a.Email
	return m, nil
}
