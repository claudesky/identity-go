package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/claudesky/identity-go/middleware"
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
	mux.HandleFunc("POST /auth/login", middleware.JSONDecoderMiddleware(c.login))
	mux.HandleFunc("GET /auth/validate", c.validate)
	mux.HandleFunc(
		"POST /auth/refresh",
		middleware.JSONDecoderMiddleware(c.refresh),
	)
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

func validateLoginRequest(rq *LoginRequest) error {
	if rq.Email == nil {
		return errors.New("[email] is required")
	}
	if rq.Password == nil {
		return errors.New("[password] is required")
	}
	return nil
}

func validateRefreshRequest(rq *RefreshRequest) error {
	if rq.RefreshToken == nil {
		return errors.New("[refresh_token] is required")
	}
	return nil
}

func (c *AuthController) login(
	w http.ResponseWriter,
	r *http.Request,
	rq *LoginRequest,
) {
	// Validation
	if err := validateLoginRequest(rq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get user
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
	err = bcrypt.CompareHashAndPassword(
		[]byte(*user.Password),
		[]byte(*rq.Password),
	)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Token Family ID
	jtf := utils.PseudoUUID()
	now := time.Now()

	// Generate Tokens
	refreshString, tokenString, expRT, err := c.th.GenerateTokens(
		&now,
		user.Id,
		jtf,
		jtf,
	)
	if err != nil {
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

func (c *AuthController) refresh(
	w http.ResponseWriter,
	r *http.Request,
	rq *RefreshRequest,
) {
	// Validation
	if err := validateRefreshRequest(rq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify Token
	ok, claims := c.th.VerifyRefreshToken(*rq.RefreshToken)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the token family
	tf, err := c.tfr.GetTokenById(r.Context(), claims.JTF)
	if err != nil {
		slog.Warn(
			"could not find token family",
			"token",
			rq.RefreshToken,
			"sub",
			claims.SUB,
			"error",
			err,
		)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Make sure sub is correct
	if claims.SUB != tf.Sub {
		slog.Warn(
			"invalid sub for refresh token",
			"token",
			rq.RefreshToken,
			"sub",
			claims.SUB,
			"family_sub",
			tf.Sub,
		)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check last issued
	if claims.JTI != tf.LastIssued {
		slog.Warn("invalid last issued for refresh token", "token", rq.RefreshToken)
		// TODO: Add revoke here
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Possible future checks
	// ...

	// New Refresh Token ID
	jtiRT := utils.PseudoUUID()

	now := time.Now()

	refreshString, tokenString, expRT, err := c.th.
		GenerateTokens(&now, tf.Sub, tf.Id, jtiRT)
	if err != nil {
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
