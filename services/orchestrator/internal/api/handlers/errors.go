package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func WriteError(w http.ResponseWriter, statusCode int, message string) {
	resp := ErrorResponse{
		Message: message,
		Code:    statusCode,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("[ERROR] Failed to write error response: %v", err)
	}
}