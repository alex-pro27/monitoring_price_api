package services

import (
	"fmt"
	"github.com/alex-pro27/monitoring_price_api/config"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"net/http"
	"net/url"
	"regexp"
)

func SendSms(phone, text string) error {
	pattern := regexp.MustCompile("(^\\+[78])|[^\\d]*")
	client := &http.Client{}
	res, err := client.PostForm(config.Config.SMS.URL, url.Values{
		"Login": {config.Config.SMS.Login},
		"Password": {config.Config.SMS.Password},
		"Source": {config.Config.SMS.Source},
		"Phone": {pattern.ReplaceAllString(phone, "")},
		"Text": {text},
	})

	if err != nil || res.StatusCode != http.StatusOK {
		logger.Logger.Errorf("Ошибка отправки смс, на номер %s", phone)
		return fmt.Errorf("Ошибка отправки смс")
	}
	logger.Logger.Infof("Отправка смс на номер %s", phone)
	return nil
}
