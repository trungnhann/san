package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type EmailSender interface {
	SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error
}

type SmtpSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
	host              string
	port              int
	username          string
}

func NewSmtpSender(name string, fromEmailAddress string, fromEmailPassword string, host string, port int, username string) EmailSender {
	if host == "" {
		host = "smtp.gmail.com"
	}
	if port == 0 {
		port = 587
	}
	if username == "" {
		username = fromEmailAddress
	}

	return &SmtpSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
		host:              host,
		port:              port,
		username:          username,
	}
}

func (sender *SmtpSender) SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.username, sender.fromEmailPassword, sender.host)
	smtpServerAddress := fmt.Sprintf("%s:%d", sender.host, sender.port)
	return e.Send(smtpServerAddress, smtpAuth)
}
