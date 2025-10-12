package controllers

import (
	"net/http"

	"github.com/claudesky/identity-go/responses"
)

func respondWithMessage(w http.ResponseWriter, message string, statusCode int) {
	responses.Message{
		Message: message,
		Status:  statusCode,
	}.Write(w)
}

func respondWithData(w http.ResponseWriter, message string, data any, statusCode int) {
	responses.DataMessage{
		Message: message,
		Data:    data,
		Status:  statusCode,
	}.Write(w)
}

// -- Preset Responses

func internalServerError(w http.ResponseWriter) {
	responses.Message{
		Status:  http.StatusInternalServerError,
		Message: "Internal Server Error",
	}.Write(w)
}

func unauthorized(w http.ResponseWriter) {
	responses.Message{
		Status:  http.StatusUnauthorized,
		Message: "Unauthorized",
	}.Write(w)
}
