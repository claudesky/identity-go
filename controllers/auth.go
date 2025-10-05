package controllers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/claudesky/identity-go/interfaces/repositories"
	"github.com/claudesky/identity-go/middleware"
	"github.com/claudesky/identity-go/models"
	"github.com/claudesky/identity-go/services"
	"github.com/claudesky/identity-go/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	th  *services.TokenHandler
	evs *services.EmailVerification
	ur  repositories.UserRepository
	tfr repositories.TokenFamilyRepository
	rrr repositories.RegisterRequestRepository
}

func NewAuthController(
	th *services.TokenHandler,
	evs *services.EmailVerification,
	ur repositories.UserRepository,
	tfr repositories.TokenFamilyRepository,
	rrr repositories.RegisterRequestRepository,
) *AuthController {
	return &AuthController{th, evs, ur, tfr, rrr}
}

func (c *AuthController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		"POST /auth/verify",
		middleware.JSONDecoderMiddleware(c.verify),
	)
	mux.HandleFunc(
		"POST /auth/register",
		middleware.JSONDecoderMiddleware(c.register),
	)
	mux.HandleFunc(
		"POST /auth/login",
		middleware.JSONDecoderMiddleware(c.login),
	)
	mux.HandleFunc(
		"GET /auth/validate",
		c.validate,
	)
	mux.HandleFunc(
		"POST /auth/refresh",
		middleware.JSONDecoderMiddleware(c.refresh),
	)
}

// -- Controller Methods

func (c *AuthController) register(
	w http.ResponseWriter,
	r *http.Request,
	rq *RegisterRequest,
) {
	// Validation
	if err := rq.validate(); err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check existing user
	if user, _ := c.ur.GetUserByEmail(r.Context(), *rq.Email); user != nil {
		respondWithError(w, "User with email already exists", http.StatusConflict)
		return
	}

	// Check existing registration request
	if rr, _ := c.rrr.GetRegisterRequestByEmail(
		r.Context(),
		*rq.Email,
	); rr != nil {
		respondWithError(
			w,
			"Registration request with email already exists",
			http.StatusConflict,
		)
		return
	}

	hashedPassword, err := utils.PasswordHash(*rq.Password)

	// Hash password
	if err != nil {
		slog.Error(
			"failed to hash password",
			"error",
			err.Error(),
			"email",
			rq.Email,
		)
		internalServerError(w)
	}

	// Create Registration Request
	registerRequest := &models.RegisterRequest{
		Id:        utils.PseudoUUID(),
		Email:     rq.Email,
		Password:  &hashedPassword,
		CreatedAt: time.Now(),
	}
	err = c.rrr.InsertRegisterRequest(r.Context(), registerRequest)
	if err != nil {
		slog.Error(
			"failed to insert registration request",
			"error",
			err.Error(),
			"email",
			rq.Email,
		)
		internalServerError(w)
		return
	}

	// Create Email Verification Request
	evr, err := c.evs.CreateEmailVerificationRequest(
		r.Context(),
		registerRequest.Id,
		*registerRequest.Email,
	)
	if err != nil {
		slog.Error(
			"failed to send email verification",
			"error",
			err.Error(),
			"email",
			rq.Email,
		)
		internalServerError(w)
		return
	}

	respondWithDataMessage(
		w,
		"Registration Request received.",
		http.StatusCreated,
		&struct {
			EmailVerificationRequest *models.EmailVerificationRequest `json:"email_verification_request"`
		}{
			EmailVerificationRequest: evr,
		},
	)
}

func (c *AuthController) verify(
	w http.ResponseWriter,
	r *http.Request,
	rq *VerificationRequest,
) {
	// Validation
	if err := rq.validate(); err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify the token
	msg, ok, err := c.evs.
		VerifyEmailVerificationRequest(r.Context(), *rq.Id, *rq.Token)
	if err != nil {
		// insert better slog error logging here
		internalServerError(w)
		return
	} else if !ok {
		respondWithError(w, msg, http.StatusUnprocessableEntity)
		return
	}

	respondWithMessage(w, "verification success", http.StatusOK)
}

func (c *AuthController) login(
	w http.ResponseWriter,
	r *http.Request,
	rq *LoginRequest,
) {
	// Validation
	if err := rq.validate(); err != nil {
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
		unauthorized(w)
		return
	}

	if user.Password == nil {
		// Better handling later
		slog.Warn("user has no password")
		unauthorized(w)
		return
	}

	// Check Password
	err = bcrypt.CompareHashAndPassword(
		[]byte(*user.Password),
		[]byte(*rq.Password),
	)
	if err != nil {
		unauthorized(w)
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
		internalServerError(w)
		return
	}

	c.tfr.InsertToken(r.Context(), &models.TokenFamily{
		Id:           jtf,
		Sub:          user.Id,
		LastIssued:   jtf,
		CreatedAt:    now,
		LastIssuedAt: now,
		ExpiresAt:    expRT,
	})

	respondWithJSON(w, &TokenResponse{
		AccessToken:  tokenString,
		RefreshToken: refreshString,
	}, http.StatusOK)
}

func (c *AuthController) refresh(
	w http.ResponseWriter,
	r *http.Request,
	rq *RefreshRequest,
) {
	// Validation
	if err := rq.validate(); err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify Token
	ok, claims := c.th.VerifyRefreshToken(*rq.RefreshToken)
	if !ok {
		unauthorized(w)
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
		unauthorized(w)
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
		unauthorized(w)
		return
	}

	// Check last issued
	if claims.JTI != tf.LastIssued {
		slog.Warn("invalid last issued for refresh token", "token", rq.RefreshToken)
		// TODO: Add revoke here
		unauthorized(w)
		return
	}

	// Possible future checks
	// ...

	// New Refresh Token ID
	jtiRT := utils.PseudoUUID()

	now := time.Now()

	// Generate new tokens
	refreshString, tokenString, expRT, err := c.th.
		GenerateTokens(&now, tf.Sub, tf.Id, jtiRT)
	if err != nil {
		internalServerError(w)
		return
	}

	// Update Token Family
	tf.LastIssued = jtiRT
	tf.LastIssuedAt = now
	tf.ExpiresAt = expRT

	c.tfr.UpdateToken(r.Context(), tf)

	// Send new tokens
	respondWithJSON(w, &TokenResponse{
		AccessToken:  tokenString,
		RefreshToken: refreshString,
	}, http.StatusOK)
}

func (c *AuthController) validate(w http.ResponseWriter, r *http.Request) {
	tokenString := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	token, err := c.th.VerifyToken(tokenString)
	if err != nil {
		slog.Info("token verification failed", "error", err)
		unauthorized(w)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		respondWithJSON(w, claims, http.StatusOK)
	} else {
		fmt.Println(err)
		internalServerError(w)
		return
	}
}
