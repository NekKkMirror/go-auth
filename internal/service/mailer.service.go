package service

import (
	"fmt"
	"github.com/NekKkMirror/go-auth/config"
	"net/smtp"
	"strings"
)

// Mailer describes the interface for sending emails.
type Mailer interface {
	Send(to, subject, body string) error
}

// SMTPMailer is responsible for sending emails via SMTP.
type SMTPMailer struct {
	cfg *config.SMTPConfig
}

// NewSMTPMailer creates an instance of SMTPMailer with the specified configuration.
func NewSMTPMailer(cfg *config.SMTPConfig) *SMTPMailer {
	return &SMTPMailer{cfg: cfg}
}

// Send sends email via SMTP.
func (m *SMTPMailer) Send(to, subject, body string) error {
	headers := map[string]string{
		"From":         m.cfg.From,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/plain; charset=\"utf-8\"",
	}
	message := ""

	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	toAddresses := strings.Split(to, ",")
	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)

	return smtp.SendMail(m.cfg.Host+":"+m.cfg.Port, auth, m.cfg.From, toAddresses, []byte(message))
}
