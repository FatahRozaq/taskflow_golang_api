package workers

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func InitFirebase() *firebase.App {
	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIALS"))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic("❌ Failed to initialize Firebase: " + err.Error())
	}
	return app
}

func SaveWeatherToFirebase(data map[string]interface{}) {
	app := InitFirebase()
	ctx := context.Background()

	client, err := app.Firestore(ctx)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	_, _, err = client.Collection("weather").Add(ctx, data)
	if err != nil {
		fmt.Println("❌ Failed to save to Firebase:", err)
	}
	fmt.Println("✅ Weather data saved to Firebase")
}
