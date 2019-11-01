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
	admins := make([]string, 0, len(config.Config.Admins))
	for _, admin := range config.Config.Admins {
		admins = append(admins, admin.Email)
	}
	SendMail(admins, subject, body, attach...)
}
