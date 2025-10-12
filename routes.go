package main

import (
	"net/http"

	"github.com/claudesky/identity-go/controllers"
	"github.com/claudesky/identity-go/middleware"
)

func registerRoutes(
	mux *http.ServeMux,
	healthController *controllers.HealthController,
	registerController *controllers.RegisterController,
	authController *controllers.AuthController,
	selfController *controllers.SelfController,
) {
	// Health routes
	mux.HandleFunc("GET /health/check", healthController.Check)

	// Register routes
	mux.HandleFunc(
		"POST /register",
		middleware.JSONDecoderMiddleware(registerController.Register),
	)
	mux.HandleFunc(
		"POST /register/verify",
		middleware.JSONDecoderMiddleware(registerController.Verify),
	)
	mux.HandleFunc(
		"POST /register/resend",
		middleware.JSONDecoderMiddleware(registerController.Resend),
	)

	// Auth routes
	mux.HandleFunc(
		"POST /auth/login",
		middleware.JSONDecoderMiddleware(authController.Login),
	)
	mux.HandleFunc(
		"GET /auth/validate",
		authController.Validate,
	)
	mux.HandleFunc(
		"POST /auth/refresh",
		middleware.JSONDecoderMiddleware(authController.Refresh),
	)

	// Self routes
	mux.HandleFunc("GET /self/sessions", selfController.Sessions)

	// Fallback Route
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})
}
