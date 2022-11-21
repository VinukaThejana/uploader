package main

import (
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func initFirebase() (*firebase.App, error) {
	config := &firebase.Config{ProjectID: os.Getenv("PROJECT_ID")}
	app, err := firebase.NewApp(ctx, config, option.WithCredentialsFile("./firebase.json"))
	return app, err
}

// VerifyIDToken - Verify the firebase idToken
func VerifyIDToken(idToken string, w http.ResponseWriter) (uid *string, status Status) {
	app, err := initFirebase()
	if err != nil {
		log.Println(err)
		return nil, InternalServerError
	}

	auth, err := app.Auth(ctx)
	if err != nil {
		log.Println(err)
		return nil, InternalServerError
	}

	token, err := auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Println(err)
		return nil, UnAuthorized
	}

	return &token.UID, Okay
}
