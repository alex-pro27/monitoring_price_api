package utils

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"github.com/alex-pro27/monitoring_price_api/config"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"github.com/go-mail/mail"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GenerateHash() string {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(strconv.FormatInt(time.Now().Unix(), 10)),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Fatal(err)
	}
	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetIPAddress(r *http.Request) string {
	var ipAddress string
	ipAddress = strings.Split(r.RemoteAddr, ":")[0]
	for _, h := range []string{"X-Forwarded-For", "X-Real-IP"} {
		for _, ip := range strings.Split(r.Header.Get(h), ",") {
			// header can contain spaces too, strip those out.
			ip = strings.TrimSpace(ip)
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() {
				// bad address, go to next
				continue
			} else {
				ipAddress = ip
				goto Done
			}
		}
	}
Done:
	return ipAddress
}

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
