package main

import (
	"net/http"

	"github.com/claudesky/identity-go/controllers"
	"github.com/claudesky/identity-go/middleware"
	"github.com/claudesky/identity-go/responses"
	"github.com/claudesky/identity-go/services"
)

func registerRoutes(
	mux *http.ServeMux,
	healthController *controllers.HealthController,
	registerController *controllers.RegisterController,
	authController *controllers.AuthController,
	selfController *controllers.SelfController,
	userController *controllers.UserController,
	tokenHandler *services.TokenHandler,
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
		"POST /auth/refresh",
		middleware.JSONDecoderMiddleware(authController.Refresh),
	)

	protected := middleware.NewGroup().
		Use(middleware.AuthMiddleware(tokenHandler)).
		Route("GET /users", userController.Index).
		Route("GET /users/{id}", userController.Show).
		Route("GET /self/sessions", selfController.Sessions).
		Route("GET /self", selfController.Self).
		Route("GET /auth/validate", authController.Validate)

	protected.Handle(mux)

	// Fallback Route
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		responses.Message{
			Message: "Not Found",
			Status:  http.StatusNotFound,
		}.Write(w)
	})
}
