package services

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App
var FcmClient *messaging.Client

func InitFirebase() {
	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIALS_PATH"))
	var err error

	FirebaseApp, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v\n", err)
	}

	FcmClient, err = FirebaseApp.Messaging(context.Background())
	if err != nil {
		log.Fatalf("error getting Firebase Messaging client: %v\n", err)
	}

	log.Println("Firebase Admin SDK initialized successfully")
}
