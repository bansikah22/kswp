package notifications

import (
	"fmt"
	"net/smtp"
)

type EmailNotifier struct {
	From     string
	To       []string
	Password string
	Host     string
	Port     string
}

func (e *EmailNotifier) Send(message string) error {
	auth := smtp.PlainAuth("", e.From, e.Password, e.Host)
	msg := []byte("To: " + e.To[0] + "\r\n" +
		"Subject: kswp notification\r\n" +
		"\r\n" +
		message + "\r\n")

	err := smtp.SendMail(e.Host+":"+e.Port, auth, e.From, e.To, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
