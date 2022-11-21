package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func initStorage() (*storage.BucketHandle, error) {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("./firebase.json"))
	if err != nil {
		return nil, err
	}

	bucket := client.Bucket(os.Getenv("STORAGE_BUCKET"))
	return bucket, err
}

// UploadFile - Upload the given file to google cloud storage bucket and return
// the URL
func UploadFile(fileName string, uid string, ext string, w http.ResponseWriter) (url *string, status Status) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Println(err.Error())
		return nil, InternalServerError
	}
	defer file.Close()

	bucket, err := initStorage()
	if err != nil {
		log.Println(err.Error())
		return nil, InternalServerError
	}
	object := bucket.Object(fmt.Sprintf("%s.%s", uid, ext))
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	wc := object.NewWriter(ctx)
	if _, err = io.Copy(wc, file); err != nil {
		log.Println(err.Error())
		return nil, InternalServerError
	}

	if err := wc.Close(); err != nil {
		log.Println(err)
		return nil, InternalServerError
	}

  URL := fmt.Sprintf("https://storage.googleapis.com/%s/%s.%s", os.Getenv("STORAGE_BUCKET"), uid, ext)
  return &URL, Okay
}
