package middleware

import (
	"log/slog"
	"net/http"

	"github.com/claudesky/identity-go/interfaces/database"
	"github.com/claudesky/identity-go/responses"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func TransactionMiddleware(db database.Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only start transaction for POST requests
			if r.Method != http.MethodPost {
				next.ServeHTTP(w, r)
				return
			}

			// Begin transaction
			txCtx, err := db.BeginTransaction(r.Context())
			if err != nil {
				slog.Error("Failed to start transaction", "error", err.Error())
				responses.Message{
					Message: "Internal Server Error",
					Status:  http.StatusInternalServerError,
				}.Write(
					w,
				)
				return
			}

			// Wrap response writer to capture status code
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK, // Default to 200
			}

			// Defer transaction handling
			defer func() {
				if rw.statusCode >= 200 && rw.statusCode < 500 {
					if err := db.CommitTransaction(txCtx); err != nil {
						slog.Error("Failed to commit transaction", "error", err.Error())
					}
				} else {
					db.RollbackTransaction(txCtx)
				}
			}()

			next.ServeHTTP(rw, r.WithContext(txCtx))
		})
	}
}
