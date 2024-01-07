package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	// smtpAuthAddress   = "smtp.qq.com"
	// smtpServerAddress = "smtp.qq.com:25"
	smtpAuthAddress   = "smtp.yeah.net"
	smtpServerAddress = "smtp.yeah.net:25"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type MailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewMailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &MailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *MailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
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

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)

	return e.Send(smtpServerAddress, smtpAuth)

}
