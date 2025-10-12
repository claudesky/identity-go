package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/claudesky/identity-go/responses"
	"github.com/claudesky/identity-go/services"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const tokenContextKey contextKey = "token"
const subjectContextKey contextKey = "subject"

func GetSubject(r *http.Request) (sub string) {
	return r.Context().Value(subjectContextKey).(string)
}

func GetToken(r *http.Request) (token *jwt.Token) {
	return r.Context().Value(tokenContextKey).(*jwt.Token)
}

func AuthMiddleware(th *services.TokenHandler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				slog.Info("missing authorization header")
				unauthorized(w)
				return
			}

			parts := strings.Split(authHeader, "Bearer ")
			if len(parts) != 2 {
				slog.Info("invalid authorization header format")
				unauthorized(w)
				return
			}

			tokenString := parts[1]
			token, err := th.VerifyToken(tokenString)
			if err != nil {
				slog.Info("token verification failed", "error", err)
				unauthorized(w)
				return
			}

			sub, err := token.Claims.GetSubject()
			if err != nil {
				slog.Warn("failed to get subject claim", "error", err)
				unauthorized(w)
				return
			}

			ctx := context.WithValue(r.Context(), tokenContextKey, token)
			ctx = context.WithValue(ctx, subjectContextKey, sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func unauthorized(w http.ResponseWriter) {
	responses.Message{
		Status:  http.StatusUnauthorized,
		Message: "Unauthorized",
	}.Write(w)
}
