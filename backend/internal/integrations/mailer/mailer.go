package mailer

import (
	"bytes"
	"fmt"
	"mime"
	"net/smtp"

	"github.com/livestreamify/backend/internal/config"
)

// Mailer sends transactional emails.
type Mailer interface {
	Send(to, subject, htmlBody string) error
}

// SMTPMailer delivers email via SMTP using net/smtp.
type SMTPMailer struct {
	host     string
	port     int
	username string
	password string
	from     string
	fromName string
}

// NewSMTPMailer creates an SMTPMailer wired from config.
func NewSMTPMailer(cfg *config.Config) *SMTPMailer {
	return &SMTPMailer{
		host:     cfg.SMTPHost,
		port:     cfg.SMTPPort,
		username: cfg.SMTPUsername,
		password: cfg.SMTPPassword,
		from:     cfg.SMTPFrom,
		fromName: cfg.SMTPFromName,
	}
}

// Send delivers an HTML email to the given address.
func (m *SMTPMailer) Send(to, subject, htmlBody string) error {
	addr := fmt.Sprintf("%s:%d", m.host, m.port)
	auth := smtp.PlainAuth("", m.username, m.password, m.host)

	fromHeader := mime.QEncoding.Encode("utf-8", m.fromName) + " <" + m.from + ">"

	var buf bytes.Buffer
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	buf.WriteString("From: " + fromHeader + "\r\n")
	buf.WriteString("To: " + to + "\r\n")
	buf.WriteString("Subject: " + mime.QEncoding.Encode("utf-8", subject) + "\r\n")
	buf.WriteString("\r\n")
	buf.WriteString(htmlBody)

	return smtp.SendMail(addr, auth, m.from, []string{to}, buf.Bytes())
}
