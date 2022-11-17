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
func VerifyIDToken(idToken string, w http.ResponseWriter) (uid string) {
	app, err := initFirebase()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	auth, err := app.Auth(ctx)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	token, err := auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	return token.UID
}
