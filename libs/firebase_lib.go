package libs

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var firebaseApp *firebase.App

func ConnectFirebase() {
	opt := option.WithCredentialsFile("firebaseServiceAccount.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Printf("error init firebase: %v", err)
	}

	firebaseApp = app
}
