package repositories

import (
	"context"

	"github.com/claudesky/identity-go/models"
)

type AuthorizationCodeRepository interface {
	GetAuthorizationCodeByCode(ctx context.Context, code string) (*models.AuthorizationCode, error)
	InsertAuthorizationCode(ctx context.Context, m *models.AuthorizationCode) error
	MarkCodeAsUsed(ctx context.Context, code string) error
}
