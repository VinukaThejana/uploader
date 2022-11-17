package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/go-redis/redis/v9"
	"google.golang.org/api/option"
)

func initVision() (*vision.ImageAnnotatorClient, error) {
	client, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsFile("./google.json"))
	return client, err
}

type response struct {
	State string `json:"state"`
}

func returnResponse(state string, sum string, redisClient *redis.Client, w http.ResponseWriter) {
	// Update the Redis database regarding the image status
	// For faster acsess
	redisClient.Set(ctx, sum, state, 0)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(&response{
		State: state,
	})
	return
}

// CheckImage - Check the image for inappropriate content
func CheckImage(fileName string, w http.ResponseWriter) {
	visionClient, err := initVision()
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer visionClient.Close()

	checkSumFile, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer checkSumFile.Close()

	// Get the checksum of the file
	hash := sha256.New()
	if _, err := io.Copy(hash, checkSumFile); err != nil {
		log.Fatal(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	sum := hex.EncodeToString(hash.Sum(nil))

	// Check the redis database for the content type of the image
	// with the check file of the image for faster processing
	redisClient, err := Redis()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	state := redisClient.Get(ctx, sum).Val()

	if state != "" {
		if state != "PROPER_CONTENT" {
			returnResponse(state, sum, redisClient, w)
		}
	}

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	props, err := visionClient.DetectSafeSearch(ctx, image, nil)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	adult := props.Adult.Enum().String()
	spoof := props.Spoof.Enum().String()
	medical := props.Medical.Enum().String()
	violence := props.Violence.Enum().String()
	racy := props.Racy.Enum().String()

	if adult == "VERY_LIKELY" || adult == "LIKELY" || adult == "POSSIBLE" {
		returnResponse("ADULT_CONTENT", sum, redisClient, w)
	}
	if spoof == "VERY_LIKELY" || spoof == "LIKELY" {
		returnResponse("SPOOF_CONTENT", sum, redisClient, w)
	}
	if medical == "VERY_LIKELY" || medical == "LIKELY" {
		returnResponse("MEDICAL_CONTENT", sum, redisClient, w)
	}
	if violence == "VERY_LIKELY" || violence == "LIKELY" {
		returnResponse("VIOLENCE_CONTENT", sum, redisClient, w)
	}
	if racy == "VERY_LIKELY" || racy == "LIKELY" {
		returnResponse("RACY_CONTENT", sum, redisClient, w)
	}

	redisClient.Set(ctx, sum, "PROPER_CONTENT", 0)
}
