package repositories

import (
	"context"

	"github.com/claudesky/identity-go/models"
)

type TokenFamilyRepository interface {
	GetTokensBySub(ctx context.Context, sub string) (*[]models.TokenFamily, error)
	InsertToken(ctx context.Context, m *models.TokenFamily) error
	GetTokenById(ctx context.Context, id string) (*models.TokenFamily, error)
	UpdateToken(ctx context.Context, m *models.TokenFamily) error
}
