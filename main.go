// Upload the given image in the multipart form
// to the sotrage buckert and CDN after testing for safe search capabilities
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

// maxUploadSize - The max upload size allowed
const maxUploadSize = 1024 * 1024 * 2 // 2MB

var ctx = context.Background()

func main() {
	// Load the env file
	if err := godotenv.Load(); err != nil {
		log.Println(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// Set the CORS policiy to the main request
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Remove all the formData from the form after execution completion
		defer func() {
			if err := r.MultipartForm.RemoveAll(); err != nil {
				log.Println(err.Error())
				HandleError(InternalServerError, w)
				return
			}
		}()

		// Limit the upload size to the specificed upload size
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			log.Println(err.Error())
			HandleError(BadRequest, w)
			return
		}

		idToken := r.FormValue("idToken")
		uid, status := VerifyIDToken(idToken, w)
		if status != Okay {
			HandleError(status, w)
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			log.Println(err.Error())
			HandleError(InternalServerError, w)
			return
		}
		defer file.Close()

		ext, status := FileTypeValidation(file, w)
		if status != Okay {
			HandleError(status, w)
		}

		// Generate a unique file name for the uploaded file
		fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
		dst, err := os.Create(fileName)
		if err != nil {
			log.Println(err)
			HandleError(InternalServerError, w)
			return
		}

		_, err = io.Copy(dst, file)
		if err != nil {
			log.Println(err)
			HandleError(InternalServerError, w)
			return
		}

		// Remove the file after completion to save
		// storage space
		defer os.Remove(fileName)

		// Check the uploaded image file for explicit content
		CheckImage(fileName, w)

		// Upload the image to the Google Cloud Storage and get the URL
		url, status := UploadFile(fileName, *uid, *ext, w)
		if status != Okay {
			if status == Json {
				return
			}

			HandleError(status, w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"URL": *url,
		})
		return
	})

	fmt.Println("Listening on port 4500")
	if err := http.ListenAndServe(":4500", mux); err != nil {
		log.Println(err.Error())
	}
}
