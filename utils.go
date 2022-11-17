package main

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

// FileTypeValidation - Detec and validate the filetype
func FileTypeValidation(file multipart.File, w http.ResponseWriter) (ext string) {
	// Creat a buffer to store thr first 512 bytes of the FileTypeValidation(
	buff := make([]byte, 512)
	_, err := file.Read(buff)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fileType := http.DetectContentType(buff)
	if fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/gif" {
		http.Error(w, "Invalid file types", http.StatusBadRequest)
		return
	}

	// Return the file pointer to the beging of the file
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get the file extension
	switch fileType {
	case "image/jpeg":
		ext = "jpg"
	case "image/png":
		ext = "png"
	case "image/gif":
		ext = "gif"
	default:
		http.Error(w, "Inavlid file type", http.StatusBadRequest)
		return
	}

	return ext
}
