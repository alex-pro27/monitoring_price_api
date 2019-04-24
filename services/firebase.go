package services

import (
	"context"
	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/alex-pro27/monitoring_price_api/config"
	"github.com/alex-pro27/monitoring_price_api/logger"
	"google.golang.org/api/option"
)

func FirebaseApp() (app *firebase.App, err error) {
	opt := option.WithCredentialsFile(config.Config.Firebase.CertPath)
	app, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logger.Logger.Error(err)
	} else {
		_, err := app.Auth(context.Background())
		logger.HandleError(err)
	}
	return app, err
}

func FirebaseSendNotification(data map[string]string, token, topic string) (err error) {
	app, err := FirebaseApp()
	if err == nil {
		messageClient, err := app.Messaging(context.Background())
		if err == nil {
			message := messaging.Message{
				Data: data,
			}
			if token != "" {
				message.Token = token
			} else {
				message.Topic = topic
			}
			_, err = messageClient.Send(context.Background(), &message)
		}
	}
	return err
}
