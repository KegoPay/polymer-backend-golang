package firebase_push_notification

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	fcm "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
	"kego.com/infrastructure/logger"
)

type FireBasePushNotification struct {
	MessagingClient *messaging.Client
}

func (fbpn *FireBasePushNotification) InitialiseClient() *FireBasePushNotification {
	firebaseServicekey := os.Getenv("FIREBASE_SERVICE_KEY")
	decodedKey, err :=  base64.StdEncoding.DecodeString(firebaseServicekey)
	if err != nil {
		logger.Error(errors.New("error converting firebase service key from base64 to byte array"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil
	}
	appInstance, err := fcm.NewApp(context.Background(), nil, option.WithCredentialsJSON(decodedKey))
	if err != nil {
		logger.Error(errors.New("error creating new app instance for firebase"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil
	}
	fbpn.MessagingClient, err = appInstance.Messaging(context.Background())
	if err != nil {
		logger.Error(errors.New("error instanciating firebase messaging client"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		})
		return nil
	}
	return fbpn
}


func (fbpn *FireBasePushNotification) PushOne(deviceID string, title string, body string) {
	_, err := fbpn.MessagingClient.Send(context.Background(), &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body: body,
		},
		Token: deviceID,
	})
	if err != nil {
		logger.Error(errors.New("error sending push notification using Firebase"), logger.LoggerOptions{
			Key: "error",
			Data: err,
		}, logger.LoggerOptions{
			Key: "deviceID",
			Data: deviceID,
		})
		return
	}
	logger.Info(fmt.Sprintf("successfully sent push notification to %s using Firebase", deviceID))
}