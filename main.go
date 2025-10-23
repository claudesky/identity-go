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
	"github.com/claudesky/identity-go/repositories/postgres"
	"github.com/claudesky/identity-go/services"
)

func main() {
	// Configure timezone to UTC
	time.Local = time.UTC

	// Configure structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Init Database
	db, err := database.NewPostgresDatabase(
		context.Background(),
		idg_db_conn,
		&idg_db_pass,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Init Repositories
	userRepository := postgres.NewUserRepository(db)
	tokenFamilyRepository := postgres.NewTokenFamilyRepository(db)
	registerRequestRepository := postgres.NewRegisterRequestRepository(db)
	emailVerificationRepository := postgres.NewEmailVerificationRepository(db)
	clientRepository := postgres.NewClientRepository(db)

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
		registerRequestRepository,
		userRepository,
		mail,
	)

	// Init Controllers
	mux := http.NewServeMux()

	healthController := controllers.NewHealthController()

	registerController := controllers.NewRegisterController(
		emailVerificationService,
		userRepository,
		registerRequestRepository,
	)

	authController := controllers.NewAuthController(
		tokenHandler,
		userRepository,
		tokenFamilyRepository,
	)

	selfController := controllers.NewSelfController(
		tokenHandler,
		userRepository,
		tokenFamilyRepository,
	)

	userController := controllers.NewUserController(
		userRepository,
	)

	clientController := controllers.NewClientController(
		clientRepository,
	)

	// Register Routes
	registerRoutes(
		mux,
		healthController,
		registerController,
		authController,
		selfController,
		userController,
		clientController,
		tokenHandler,
	)

	// Start Server

	slog.Info("server init")
	log.Fatal(
		http.ListenAndServe(
			idg_port,
			middleware.TransactionMiddleware(db)(
				middleware.ContentTypeJson(mux),
			),
		),
	)
}
