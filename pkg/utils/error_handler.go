package utils

import (
	"encoding/json"
	"net/http"
)

// Error struct will be used to represent errors as JSON
type Error struct {
	Message string `json:"message"`
}

// WriteError writes an error message as JSON to a http.ResponseWriter
func WriteError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	errorResponse := Error{Message: message}
	json.NewEncoder(w).Encode(errorResponse)
}
