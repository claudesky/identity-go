package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/claudesky/identity-go/repositories"
	"github.com/claudesky/identity-go/services"
	"github.com/claudesky/identity-go/utils"
	"github.com/golang-jwt/jwt/v5"
)

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type AuthController struct {
	tokenHandler   *services.TokenHandler
	userRepository *repositories.UserRepository
}

func NewAuthController(
	th *services.TokenHandler,
	ur *repositories.UserRepository,
) *AuthController {
	return &AuthController{
		tokenHandler:   th,
		userRepository: ur,
	}
}

func (c *AuthController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/login", c.login)
	mux.HandleFunc("GET /auth/validate", c.validate)
}

func (c *AuthController) login(w http.ResponseWriter, r *http.Request) {
	email := "admin@example.org"
	user, err := c.userRepository.GetUserByEmail(r.Context(), email)
	if err != nil {
		slog.Info("could not find user by email",
			slog.String("error", fmt.Sprintf("%v", err)),
			slog.Group("params",
				slog.String("email", email),
			),
		)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	// Assume Login Success

	// Token Family ID
	jtf := utils.PseudoUUID()

	// Tokens TTL
	ttlRT := time.Hour * time.Duration(72)
	ttlAT := time.Minute * time.Duration(5)

	refreshString, err := c.tokenHandler.SignToken(jwt.MapClaims{
		"jti": jtf,
		"jtf": jtf,
		"exp": time.Now().UTC().Add(ttlRT).Unix(),
	})
	if err != nil {
		slog.Info("refresh token signing failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tokenString, err := c.tokenHandler.SignToken(jwt.MapClaims{
		"jti": utils.PseudoUUID(),
		"jtf": jtf,
		"jtp": jtf,
		"sub": user.Id,
		"exp": time.Now().UTC().Add(ttlAT).Unix(),
	})
	if err != nil {
		log.Fatalf("access token signing failed: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Save the token family here

	json.NewEncoder(w).Encode(&TokenResponse{
		AccessToken:  tokenString,
		RefreshToken: refreshString,
	})
}

func (c *AuthController) validate(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	token, err := c.tokenHandler.VerifyToken(tokenString)
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
