package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/claudesky/identity-go/controllers"
	"github.com/claudesky/identity-go/database"
	"github.com/claudesky/identity-go/middleware"
	"github.com/claudesky/identity-go/repositories"
	"github.com/claudesky/identity-go/services"
)

func main() {
	// Configure timezone to UTC
	time.Local = time.UTC

	// Configure structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Init Database
	database, err := database.NewDatabase(
		context.Background(),
		idg_db_conn,
		&idg_db_pass,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Init Repositories
	userRepository := repositories.NewUserRepository(database)
	tokenFamilyRepository := repositories.NewTokenFamilyRepository(database)
	registerRequestRepository := repositories.NewRegisterRequestRepository(
		database,
	)
	emailVerificationRepository := repositories.NewEmailVerificationRepository(
		database,
	)

	// Init Services
	tokenHandler := services.NewTokenHandler(idg_pkey, idg_pubkey)
	mail := services.NewMail(
		idg_mail_user,
		idg_mail_pass,
		idg_mail_host,
		idg_mail_port,
		idg_mail_addr,
	)
	emailVerificationService := services.NewEmailVerification(
		emailVerificationRepository,
		mail,
	)

	// Init Controllers
	mux := http.NewServeMux()

	healthController := controllers.NewHealthController()
	healthController.RegisterRoutes(mux)

	authController := controllers.NewAuthController(
		tokenHandler,
		emailVerificationService,
		userRepository,
		tokenFamilyRepository,
		registerRequestRepository,
	)
	authController.RegisterRoutes(mux)

	selfController := controllers.NewSelfController(
		tokenHandler,
		userRepository,
		tokenFamilyRepository,
	)
	selfController.RegisterRoutes(mux)

	// Fallback Route
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	// Start Server

	// Send Init Email
	if idg_send_init_email {
		if err := mail.SendSystemEmailSimple("admin@example.org", "test", "server init"); err != nil {
			slog.Warn("Initialization Email Error", "error", err.Error())
		}
	}

	slog.Info("server init")
	log.Fatal(http.ListenAndServe(idg_port, middleware.ContentTypeJson(mux)))
}
