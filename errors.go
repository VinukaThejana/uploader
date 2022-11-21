package main

import "net/http"

// Status - Status codes
type Status int64

const (
	// Okay - A normal response without an InternalServerError
	Okay Status = iota
	// BadRequest - A bad request sent by the user
	BadRequest
	// InternalServerError - An Internal server error occured
	InternalServerError
	// UnAuthorized - The user is unauthorized to perform the given action
	UnAuthorized
	// MethodNotAllowed - The Given method is not allowed
	MethodNotAllowed
	// Json - The Json response
	Json
)

// HandleError - Handle the given error
func HandleError(err Status, w http.ResponseWriter) {
	switch err {
	case BadRequest:
		http.Error(w, "Bad request", http.StatusBadRequest)
	case InternalServerError:
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	case UnAuthorized:
		http.Error(w, "UnAuthorized", http.StatusUnauthorized)
	case MethodNotAllowed:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
