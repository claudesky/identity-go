package main

import (
	"crypto"
	"flag"
	"fmt"
	"log"
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

	tokenHandler := services.NewTokenHandler(pkey, pubkey)

	mux := http.NewServeMux()

	healthController := controllers.HealthController{Healthy: true}
	healthController.RegisterRoutes(mux)

	authController := controllers.AuthController{
		TokenHandler: tokenHandler,
	}
	authController.RegisterRoutes(mux)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	log.Fatal(http.ListenAndServe(*addr, mux))
}

func fetchPkey() (crypto.PrivateKey, error) {
	// Get the Private Key
	pkeyBytes, err := os.ReadFile(*privKey)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	pkey, err := jwt.ParseEdPrivateKeyFromPEM(pkeyBytes)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return pkey, nil
}

func fetchPubkey() (crypto.PublicKey, error) {
	// Get the Public Key
	pubkeyBytes, err := os.ReadFile(*pubKey)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	pubkey, err := jwt.ParseEdPublicKeyFromPEM(pubkeyBytes)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return pubkey, nil
}
