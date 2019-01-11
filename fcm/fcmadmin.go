package fcm

import (
	"context"
	"log"

	"firebase.google.com/go"
	"google.golang.org/api/option"
	"firebase.google.com/go/messaging"
	"fmt"
	"github.com/InclusION/model"
	"time"
)

func InitializeAppWithServiceAccount() *firebase.App {
	// [START initialize_app_service_account_golang]
	opt := option.WithCredentialsFile("adminsdk.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	// [END initialize_app_service_account_golang]

	return app
}


func Notify1Device(deviceToken string, mdata string) error {
	err := sendTo1Device(deviceToken, mdata)
	if err != nil {
		return err
	}
	return nil
}


func NotifyAllDevicesOfUser(username string, mdata string) error {

	// query all user's devices
	p := model.Phone{}
	err, phones := p.QueryPhonesByUsername(username)
	if err != nil {
		return err
	}

	for _, phone := range phones {
		err = sendTo1Device(phone.FCMToken, mdata)
		if err != nil {
			log.Printf("Error sendTo1Device %s", err)
		}
	}

	return nil
}

func sendTo1Device(deviceToken string, mdata string) error {

	app := InitializeAppWithServiceAccount()

	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Printf("error getting Messaging client: %v\n", err)
		return err
	}

	// See documentation on defining a message payload.
	message := &messaging.Message {
		Data: map[string]string {
			"mdata": mdata,
			"time":  time.Now().String(),
		},
		Token: deviceToken,
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, message)
	if err != nil {
		log.Println(err)
		return err
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)

	return nil
}

