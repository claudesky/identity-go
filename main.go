package main

import (
	"crypto"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/claudesky/identity-go/controllers"
	"github.com/claudesky/identity-go/services"
	"github.com/golang-jwt/jwt/v5"
)

var addr = flag.String("addr", ":9102", "Http Service Address")
var privKey = flag.String("privKey", "./keys/private.pem", "Private Key")
var pubKey = flag.String("pubKey", "./keys/public.pem", "Public Key")

func main() {
	flag.Parse()

	// Configure structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	// Fetch keys
	pkey, err := fetchPkey()
	if err != nil {
		log.Fatal(err)
		return
	}

	pubkey, err := fetchPubkey()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Init Services
	tokenHandler := services.NewTokenHandler(pkey, pubkey)

	// Init Controllers
	mux := http.NewServeMux()

	healthController := controllers.HealthController{Healthy: true}
	healthController.RegisterRoutes(mux)

	authController := controllers.AuthController{
		TokenHandler: tokenHandler,
	}
	authController.RegisterRoutes(mux)

	// Fallback Route
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	// Start Server
	slog.Info("server init")
	log.Fatal(http.ListenAndServe(*addr, mux))
}

func fetchPkey() (crypto.PrivateKey, error) {
	// Get the Private Key
	pkeyBytes, err := os.ReadFile(*privKey)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jwt.ParseEdPrivateKeyFromPEM(pkeyBytes)
}

func fetchPubkey() (crypto.PublicKey, error) {
	// Get the Public Key
	pubkeyBytes, err := os.ReadFile(*pubKey)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return jwt.ParseEdPublicKeyFromPEM(pubkeyBytes)
}
