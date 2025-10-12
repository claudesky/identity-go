package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/claudesky/identity-go/requests"
	"github.com/claudesky/identity-go/responses"
)

type validateablePtr[T any] interface {
	requests.Validateable
	*T
}

type HandlerWithRequest[T any] func(w http.ResponseWriter, r *http.Request, body *T)

func JSONDecoderMiddleware[T any, PT validateablePtr[T]](
	next HandlerWithRequest[T],
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body T

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			responses.Message{
				Message: "Invalid JSON",
				Status:  http.StatusBadRequest,
			}.Write(w)
			slog.Info("Received Bad Request", "error", err.Error())
			return
		}

		if err := PT(&body).Validate(); err != nil {
			responses.Message{
				Message: err.Error(),
				Status:  http.StatusUnprocessableEntity,
			}.Write(w)
			slog.Info("Received Unprocessable Entity", "error", err.Error())
			return
		}

		next(w, r, &body)
	}
}
