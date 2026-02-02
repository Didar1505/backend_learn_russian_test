package mailer

import (
	"context"
	"github.com/Didar1505/project_test.git/pkg/config"
	"gopkg.in/gomail.v2"
)

type SMTPSender struct {
	dialer *gomail.Dialer
	from   string
}

func New(cfg *config.Config) *SMTPSender {
	d := gomail.NewDialer(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUser,
		cfg.SMTPPass,
	)

	d.SSL = false

	return &SMTPSender{
		dialer: d,
		from:   cfg.SMTPFrom,
	}
}

func (s *SMTPSender) SendOTP(ctx context.Context, to, body, textBody string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Your OTP Code")
	m.SetBody("text/plain; charset=UTF-8", textBody)
	m.AddAlternative("text/html; charset=UTF-8", body)
	// log.Printf("[DEV OTP] email=%s code=%s\n", to, body)
	return s.dialer.DialAndSend(m)
}
