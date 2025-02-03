package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type HandlerWithRequest[T any] func(w http.ResponseWriter, r *http.Request, body *T)

func JSONDecoderMiddleware[T any](next HandlerWithRequest[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body T

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			slog.Warn("Received Bad Request", "error", err.Error())
			return
		}

		next(w, r, &body)
	}
}
