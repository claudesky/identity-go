package controllers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/claudesky/identity-go/interfaces/repositories"
	"github.com/claudesky/identity-go/models"
	"github.com/claudesky/identity-go/requests"
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

// -- Controller Methods

func (c *RegisterController) Register(
	w http.ResponseWriter,
	r *http.Request,
	rq *requests.Register,
) {
	// Check existing user
	if user, _ := c.ur.GetUserByEmail(r.Context(), *rq.Email); user != nil {
		respondWithMessage(w, "User with email already exists", http.StatusConflict)
		return
	}

	// Check existing registration request
	if rr, _ := c.rrr.GetRegisterRequestByEmail(
		r.Context(),
		*rq.Email,
	); rr != nil {
		respondWithMessage(
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

	respondWithData(
		w,
		"Registration Request received.",
		&struct {
			EmailVerificationRequest *models.EmailVerificationRequest `json:"email_verification_request"`
		}{
			EmailVerificationRequest: evr,
		},
		http.StatusCreated,
	)
}

func (c *RegisterController) Resend(
	w http.ResponseWriter,
	r *http.Request,
	rq *requests.ResendVerification,
) {
	// Check if user already exists
	if user, _ := c.ur.GetUserByEmail(r.Context(), *rq.Email); user != nil {
		respondWithMessage(w, "User with email already exists", http.StatusConflict)
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

	respondWithData(
		w,
		"Verification email resent.",
		&struct {
			EmailVerificationRequest *models.EmailVerificationRequest `json:"email_verification_request"`
		}{
			EmailVerificationRequest: evr,
		},
		http.StatusOK,
	)
}

func (c *RegisterController) Verify(
	w http.ResponseWriter,
	r *http.Request,
	rq *requests.Verification,
) {
	// Verify the token
	msg, ok, err := c.evs.
		VerifyEmailVerificationRequest(r.Context(), *rq.Id, *rq.Token)
	if err != nil {
		// insert better slog error logging here
		slog.Error(err.Error())
		internalServerError(w)
		return
	} else if !ok {
		respondWithMessage(w, msg, http.StatusUnprocessableEntity)
		return
	}

	respondWithMessage(w, "Success", http.StatusOK)
}
