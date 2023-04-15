package service

import (
	"context"
	"encoding/base64"
	"log"
	"os"

	"google.golang.org/api/option"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

func PushNotification(deviceToken string, title string, body string) error {
	var deviceIds []string
	deviceIds = append(deviceIds, deviceToken)
	return SendPushNotification(deviceIds, title, body)
}

func SendPushNotification(deviceTokens []string, title string, body string) error {
	decodedKey, err := getDecodedFireBaseKey()
	if err != nil {
		return err
	}

	opts := []option.ClientOption{option.WithCredentialsJSON(decodedKey)}

	app, err := firebase.NewApp(context.Background(), nil, opts...)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	fcmClient, err := app.Messaging(context.Background())
	if err != nil {
		return err
	}

	response, err := fcmClient.SendMulticast(context.Background(), &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Tokens: deviceTokens,
	})

	if err != nil {
		return err
	}

	log.Println("Response success count : ", response.SuccessCount)
	log.Println("Response failure count : ", response.FailureCount)

	return nil
}

func getDecodedFireBaseKey() ([]byte, error) {
	fireBaseAuthKey := os.Getenv("FIREBASE_AUTH_KEY")
	decodedKey, err := base64.StdEncoding.DecodeString(fireBaseAuthKey)
	if err != nil {
		return nil, err
	}

	return decodedKey, nil
}
