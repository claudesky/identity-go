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
	ev, err := c.evr.GetEmailVerificationRequestById(ctx, id)
	if err != nil {
		msg = "Token not found"
		return
	}

	// Check if invalid attempt was previously made
	// and subsequently the verification request was revoked
	if ev.Revoked {
		msg = "Token is revoked"
		return
	}

	// Check if already accepted
	if ev.Accepted {
		msg = "Token is already accepted"
		return
	}

	// Check if expired
	if time.Now().After(ev.ExpiresAt) {
		msg = "Token is expired"
		return
	}

	// Check if token matches
	if *ev.Token != token {
		msg = "Token did not match"
		// Revoke if not matching
		c.evr.Revoke(ctx, ev)
		return
	}

	// Everything OK
	ok = true

	rr, err := c.rr.GetRegisterRequestById(ctx, *ev.RegisterRequestId)
	if err != nil {
		return
	}

	// Invalidate the RRR
	// Probably also include a link to the created user?
	// Revoke all remaining email verification requests as well.
	// Can still return an error from here? If RRR was already revoked

	err = c.ur.InsertUser(ctx, models.NewUser(*rr.Password, *rr.Email))
	if err != nil {
		return
	}

	return
}
