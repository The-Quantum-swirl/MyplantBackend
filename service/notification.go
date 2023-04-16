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
	decodedKey, err := getDecodedFireBaseKey()
	if err != nil {
		return err
	}

	opts := []option.ClientOption{option.WithCredentialsJSON(decodedKey)}

	app, err := firebase.NewApp(context.TODO(), nil, opts...)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	fcmClient, err := app.Messaging(context.TODO())
	if err != nil {
		return err
	}

	response, err := fcmClient.Send(context.TODO(), &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: deviceToken,
	})

	if err != nil {
		log.Fatalf("error while pushing %e", err)
		return err
	}

	log.Printf("Response : {%s}", response)
	return nil
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
