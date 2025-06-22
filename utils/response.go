package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, message string, status int) {
	WriteJSON(w, status, map[string]string{"error": message})
}
