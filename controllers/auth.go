package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/claudesky/identity-go/services"
	"github.com/claudesky/identity-go/utils"
	"github.com/golang-jwt/jwt/v5"
)

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type AuthController struct {
	TokenHandler *services.TokenHandler
}

func (c *AuthController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/login", c.login)
	mux.HandleFunc("GET /auth/validate", c.validate)
}

func (c *AuthController) login(w http.ResponseWriter, _ *http.Request) {
	// TODO: Login Stuff here

	// Assume Login Success

	// Token Family ID
	jtf := utils.PseudoUUID()

	// Tokens TTL
	ttlRT := time.Hour * time.Duration(72)
	ttlAT := time.Minute * time.Duration(5)

	refreshString, err := c.TokenHandler.SignToken(jwt.MapClaims{
		"jti": jtf,
		"jtf": jtf,
		"exp": time.Now().UTC().Add(ttlRT).Unix(),
	})
	if err != nil {
		slog.Info("refresh token signing failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tokenString, err := c.TokenHandler.SignToken(jwt.MapClaims{
		"jti": utils.PseudoUUID(),
		"jtf": jtf,
		"jtp": jtf,
		"sub": "subject",
		"exp": time.Now().UTC().Add(ttlAT).Unix(),
	})
	if err != nil {
		log.Fatalf("access token signing failed: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&TokenResponse{
		AccessToken:  tokenString,
		RefreshToken: refreshString,
	})
}

func (c *AuthController) validate(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	token, err := c.TokenHandler.VerifyToken(tokenString)
	if err != nil {
		slog.Info("token verification failed", "error", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		json.NewEncoder(w).Encode(claims)
	} else {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
