package email

import (
	"milliard-easy/daaz_simbank/context"
	"milliard-easy/daaz_simbank/log"

	"github.com/sirupsen/logrus"
	gomail "gopkg.in/mail.v2"
)

// SendEmail message to config user
func SendEmail(c *context.Config, subject, body string) {
	l := logrus.WithFields(logrus.Fields{
		"subject": subject,
		"body":    body,
	})
	l.Debugf(log.DebugColor, "Sending email")

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", c.SMTP.From)

	// Set E-Mail receivers
	m.SetHeader("To", c.SMTP.To)

	// Set E-Mail subject
	m.SetHeader("Subject", subject)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", body)

	// Settings for SMTP server
	d := gomail.NewDialer(c.SMTP.Host, c.SMTP.Port, c.SMTP.From, c.SMTP.Password)

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		l.WithError(err).Errorf(log.ErrorColor, "Cannot send email")
	}
}
