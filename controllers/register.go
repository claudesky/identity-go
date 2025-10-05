package controllers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/claudesky/identity-go/interfaces/repositories"
	"github.com/claudesky/identity-go/middleware"
	"github.com/claudesky/identity-go/models"
	"github.com/claudesky/identity-go/services"
	"github.com/claudesky/identity-go/utils"
)

type RegisterController struct {
	evs *services.EmailVerification
	ur  repositories.UserRepository
	rrr repositories.RegisterRequestRepository
}

func NewRegisterController(
	evs *services.EmailVerification,
	ur repositories.UserRepository,
	rrr repositories.RegisterRequestRepository,
) *RegisterController {
	return &RegisterController{evs, ur, rrr}
}

func (c *RegisterController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		"POST /register",
		middleware.JSONDecoderMiddleware(c.register),
	)
	mux.HandleFunc(
		"POST /register/verify",
		middleware.JSONDecoderMiddleware(c.verify),
	)
	mux.HandleFunc(
		"POST /register/resend",
		middleware.JSONDecoderMiddleware(c.resend),
	)
}

// -- Controller Methods

func (c *RegisterController) register(
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

func (c *RegisterController) resend(
	w http.ResponseWriter,
	r *http.Request,
	rq *ResendVerificationRequest,
) {
	// Validation
	if err := rq.validate(); err != nil {
		respondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user already exists
	if user, _ := c.ur.GetUserByEmail(r.Context(), *rq.Email); user != nil {
		respondWithError(w, "User with email already exists", http.StatusConflict)
		return
	}

	// Resend verification email
	evr, err := c.evs.ResendEmailVerificationRequest(r.Context(), *rq.Email)
	if err != nil {
		slog.Error(
			"failed to resend email verification",
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
		"Verification email resent.",
		http.StatusOK,
		&struct {
			EmailVerificationRequest *models.EmailVerificationRequest `json:"email_verification_request"`
		}{
			EmailVerificationRequest: evr,
		},
	)
}

func (c *RegisterController) verify(
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
