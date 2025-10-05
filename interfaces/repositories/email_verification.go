package repositories

import (
	"context"

	"github.com/claudesky/identity-go/models"
)

type EmailVerificationRepository interface {
	InsertEmailVerificationRequest(ctx context.Context, m *models.EmailVerificationRequest) error
	GetEmailVerificationRequestsByRegisterRequestId(ctx context.Context, id string) (*[]models.EmailVerificationRequest, error)
	GetEmailVerificationRequestById(ctx context.Context, id string) (*models.EmailVerificationRequest, error)
	Revoke(ctx context.Context, m *models.EmailVerificationRequest) error
	Accept(ctx context.Context, m *models.EmailVerificationRequest) error
}
