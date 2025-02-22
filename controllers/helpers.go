package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// -- Response Generators

func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&ErrorMessage{
		Error:  message,
		Status: statusCode,
	})
}

func respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func respondWithMessage(w http.ResponseWriter, message string, statusCode int) {
	respondWithJSON(w, &Message{Message: message, Status: statusCode}, statusCode)
}

func respondWithDataMessage(
	w http.ResponseWriter,
	message string,
	statusCode int,
	data any,
) {
	respondWithJSON(w, &DataMessage{
		Message: message,
		Status:  statusCode,
		Data:    data,
	}, statusCode)
}

// -- Preset Responses

func internalServerError(w http.ResponseWriter) {
	respondWithError(w, "Internal Server Error", http.StatusInternalServerError)
}

func unauthorized(w http.ResponseWriter) {
	respondWithError(w, "Unauthorized", http.StatusUnauthorized)
}

// -- Logging Helpers

func logError(logger slog.Logger, err error, message string, context ...string) {

}
