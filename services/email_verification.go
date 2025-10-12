package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/claudesky/identity-go/interfaces/repositories"
	"github.com/claudesky/identity-go/models"
)

type EmailVerification struct {
	evr  repositories.EmailVerificationRepository
	rr   repositories.RegisterRequestRepository
	ur   repositories.UserRepository
	mail *Mail
}

func NewEmailVerification(
	evr repositories.EmailVerificationRepository,
	rr repositories.RegisterRequestRepository,
	ur repositories.UserRepository,
	mail *Mail,
) *EmailVerification {
	return &EmailVerification{evr, rr, ur, mail}
}

func (c *EmailVerification) CreateEmailVerificationRequest(
	ctx context.Context,
	rrid string,
	email string,
) (*models.EmailVerificationRequest, error) {
	evr := models.NewEmailVerificationRequest(rrid, email)

	// Send Email
	err := c.mail.SendSystemEmailSimple(
		*evr.Email,
		"Identity Verification",
		fmt.Sprintf("Your Identity verification code is: %s", *evr.Token),
	)

	if err != nil {
		return nil, err
	}

	err = c.evr.InsertEmailVerificationRequest(ctx, evr)

	if err != nil {
		slog.Error("Email Verification Request insert failed", "error", err.Error())
		return nil, err
	}

	return evr, nil
}

func (c *EmailVerification) ResendEmailVerificationRequest(
	ctx context.Context,
	email string,
) (*models.EmailVerificationRequest, error) {
	// Look up the register request by email
	rr, err := c.rr.GetRegisterRequestByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Revoke all previous valid verification requests
	if err := c.evr.RevokeAllByRegisterRequestId(ctx, rr.Id); err != nil {
		slog.Error(
			"Failed to revoke previous verification requests",
			"error",
			err.Error(),
			"register_request_id",
			rr.Id,
		)
		return nil, err
	}

	// Create and send new verification request
	return c.CreateEmailVerificationRequest(ctx, rr.Id, email)
}

func (c *EmailVerification) VerifyEmailVerificationRequest(
	ctx context.Context,
	id string,
	token string,
) (msg string, ok bool, err error) {
	evr, err := c.evr.GetEmailVerificationRequestById(ctx, id)
	if err != nil {
		msg = "Token not found"
		return
	}

	// Check if invalid attempt was previously made
	// and subsequently the verification request was revoked
	if evr.Revoked {
		msg = "Token is revoked"
		return
	}

	// Check if already accepted
	if evr.Accepted {
		// no need to error out here, response should be OK and user should
		// already be created
		ok = true
		return
	}

	// Check if expired
	if time.Now().After(evr.ExpiresAt) {
		msg = "Token is expired"
		return
	}

	// Check if token matches
	if *evr.Token != token {
		msg = "Token did not match"
		// Revoke if not matching
		c.evr.Revoke(ctx, evr)
		return
	}

	// Everything OK
	ok = true

	rr, err := c.rr.GetRegisterRequestById(ctx, *evr.RegisterRequestId)
	if err != nil {
		return
	}

	err = c.evr.Accept(ctx, evr)
	if err != nil {
		return
	}

	newUser := models.NewUser(*rr.Password, *rr.Email)
	now := time.Now()
	newUser.EmailVerifiedOn = &now

	err = c.ur.InsertUser(ctx, newUser)
	if err != nil {
		return
	}

	return
}
