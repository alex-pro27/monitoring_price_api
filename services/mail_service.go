package services

import (
	"crypto/tls"
	"github.com/alex-pro27/monitoring_price_api/config"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/go-mail/mail"
)

func SendMail(to []string, subject, body string, attach ...string) {
	go func() {
		emailConf := config.Config.Email
		d := mail.Dialer{
			Host:      emailConf.Host,
			Port:      emailConf.Port,
			Username:  emailConf.User,
			Password:  emailConf.Password,
			TLSConfig: &tls.Config{InsecureSkipVerify: true},
		}
		m := mail.NewMessage()
		m.SetAddressHeader("From", emailConf.From, emailConf.Name)
		m.SetHeader("To", to...)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", body)
		for _, path := range attach {
			m.Attach(path)
		}
		err := d.DialAndSend(m)
		logger.HandleError(err)
	}()
}

func SendMailToAdmin(subject, body string, attach ...string) {
	SendMail([]string{config.Config.Admin.Email}, subject, body, attach...)
}
