package repositories

import (
	"context"

	"github.com/claudesky/identity-go/models"
)

type RegisterRequestRepository interface {
	InsertRegisterRequest(ctx context.Context, m *models.RegisterRequest) error
	GetRegisterRequestById(ctx context.Context, id string) (*models.RegisterRequest, error)
	GetRegisterRequestByEmail(ctx context.Context, email string) (*models.RegisterRequest, error)
}
