package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/claudesky/identity-go/controllers"
	"github.com/claudesky/identity-go/repositories"
	"github.com/claudesky/identity-go/services"
)

func main() {
	// Configure structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Init Services
	tokenHandler := services.NewTokenHandler(idg_pkey, idg_pubkey)
	database, err := services.NewDatabase(
		context.Background(),
		idg_db_conn,
		&idg_db_pass,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Init Repositories
	userRepository := repositories.NewUserRepository(database)

	// Init Controllers
	mux := http.NewServeMux()

	healthController := controllers.NewHealthController()
	healthController.RegisterRoutes(mux)

	authController := controllers.NewAuthController(tokenHandler, userRepository)
	authController.RegisterRoutes(mux)

	// Fallback Route
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	// Start Server
	slog.Info("server init")
	log.Fatal(http.ListenAndServe(idg_port, mux))
}
