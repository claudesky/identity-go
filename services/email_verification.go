package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/claudesky/identity-go/models"
	"github.com/claudesky/identity-go/repositories"
)

type EmailVerification struct {
	evr  *repositories.EmailVerificationRepository
	mail *Mail
}

func NewEmailVerification(
	evr *repositories.EmailVerificationRepository,
	mail *Mail,
) *EmailVerification {
	return &EmailVerification{evr, mail}
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
