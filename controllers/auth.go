package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/claudesky/identity-go/models"
	"github.com/claudesky/identity-go/repositories"
	"github.com/claudesky/identity-go/services"
	"github.com/claudesky/identity-go/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	th  *services.TokenHandler
	ur  *repositories.UserRepository
	tfr *repositories.TokenFamilyRepository
}

func NewAuthController(
	th *services.TokenHandler,
	ur *repositories.UserRepository,
	tfr *repositories.TokenFamilyRepository,
) *AuthController {
	return &AuthController{th, ur, tfr}
}

func (c *AuthController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/login", c.login)
	mux.HandleFunc("GET /auth/validate", c.validate)
	mux.HandleFunc("POST /auth/refresh", c.refresh)
}

type LoginRequest struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken *string `json:"refresh_token"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (c *AuthController) login(w http.ResponseWriter, r *http.Request) {
	var rq *LoginRequest

	rq, err := utils.DecodeRequestJSON[LoginRequest](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validation
	if rq.Email == nil {
		http.Error(w, "[email] is required", http.StatusBadRequest)
		return
	}

	if rq.Password == nil {
		http.Error(w, "[password] is required", http.StatusBadRequest)
		return
	}

	user, err := c.ur.GetUserByEmail(r.Context(), *rq.Email)
	if err != nil {
		// Horrible error handling, should get this handled outside?
		slog.Info("could not find user by email",
			slog.String("error", fmt.Sprintf("%v", err)),
			slog.Group("params",
				slog.String("email", *rq.Email),
			),
		)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if user.Password == nil {
		// Better handling later
		slog.Warn("user has no password")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check Password
	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(*rq.Password))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Token Family ID
	jtf := utils.PseudoUUID()

	// Tokens TTL
	now := time.Now()
	ttlRT := time.Hour * time.Duration(72)
	ttlAT := time.Minute * time.Duration(5)
	expRT := now.Add(ttlRT)
	expAT := now.Add(ttlAT)

	refreshString, err := c.th.SignToken(jwt.MapClaims{
		"jti": jtf,
		"jtf": jtf,
		"sub": user.Id,
		"exp": expRT.Unix(),
		"typ": "refresh_token",
	})
	if err != nil {
		slog.Info("refresh token signing failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tokenString, err := c.th.SignToken(jwt.MapClaims{
		"jti": utils.PseudoUUID(),
		"jtf": jtf,
		"jtp": jtf,
		"sub": user.Id,
		"exp": expAT.Unix(),
		"typ": "access_token",
	})
	if err != nil {
		log.Fatalf("access token signing failed: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	c.tfr.InsertToken(r.Context(), &models.TokenFamily{
		Id:           jtf,
		Sub:          user.Id,
		LastIssued:   jtf,
		CreatedAt:    now.UTC(),
		LastIssuedAt: now.UTC(),
		ExpiresAt:    expRT.UTC(),
	})

	json.NewEncoder(w).Encode(&TokenResponse{
		AccessToken:  tokenString,
		RefreshToken: refreshString,
	})
}

func (c *AuthController) refresh(w http.ResponseWriter, r *http.Request) {
	var rq *RefreshRequest

	rq, err := utils.DecodeRequestJSON[RefreshRequest](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if rq.RefreshToken == nil {
		http.Error(w, "[refresh_token] is required", http.StatusBadRequest)
		return
	}

	// Verify Token
	token, err := c.th.VerifyToken(*rq.RefreshToken)

	if err != nil {
		slog.Info("token verification failed", "error", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Verify token is a refresh token
	jti, ok := claims["jti"].(string)
	if !ok {
		slog.Error("no jti", "token", rq.RefreshToken)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	sub, err := claims.GetSubject()
	if err != nil {
		slog.Error("no sub", "token", rq.RefreshToken, "error", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	typ, ok := claims["typ"].(string)
	if !ok {
		slog.Error("no typ", "token", rq.RefreshToken)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	jtf, ok := claims["jtf"].(string)
	if !ok {
		slog.Error("no jtf", "token", rq.RefreshToken)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if typ != "refresh_token" {
		http.Error(w, "token must be a refresh token", http.StatusBadRequest)
		return
	}

	// Get the token family
	tf, err := c.tfr.GetTokenById(r.Context(), jtf)
	if err != nil {
		slog.Warn("could not find token family", "token", rq.RefreshToken, "sub", sub, "error", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Make sure sub is correct
	if sub != tf.Sub {
		slog.Warn("invalid sub for refresh token", "token", rq.RefreshToken, "sub", sub, "family_sub", tf.Sub)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check last issued
	if jti != tf.LastIssued {
		slog.Warn("invalid last issued for refresh token", "token", rq.RefreshToken)
		// TODO: Add revoke here
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Possible future checks

	// Redundant logic with login, for refactor

	// New Refresh Token ID
	jtiRT := utils.PseudoUUID()

	// Tokens TTL
	now := time.Now()
	ttlRT := time.Hour * time.Duration(72)
	ttlAT := time.Minute * time.Duration(5)
	expRT := now.Add(ttlRT)
	expAT := now.Add(ttlAT)

	refreshString, err := c.th.SignToken(jwt.MapClaims{
		"jti": jtiRT,
		"jtf": jtf,
		"sub": sub,
		"exp": expRT.Unix(),
		"typ": "refresh_token",
	})
	if err != nil {
		slog.Info("refresh token signing failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tokenString, err := c.th.SignToken(jwt.MapClaims{
		"jti": utils.PseudoUUID(),
		"jtf": jtf,
		"jtp": jtiRT,
		"sub": sub,
		"exp": expAT.Unix(),
		"typ": "access_token",
	})
	if err != nil {
		log.Fatalf("access token signing failed: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Update Token Family
	tf.LastIssued = jtiRT
	tf.LastIssuedAt = now.UTC()
	tf.ExpiresAt = expRT.UTC()

	c.tfr.UpdateToken(r.Context(), tf)

	// Send new tokens
	json.NewEncoder(w).Encode(&TokenResponse{
		AccessToken:  tokenString,
		RefreshToken: refreshString,
	})
}

func (c *AuthController) validate(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	token, err := c.th.VerifyToken(tokenString)
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
