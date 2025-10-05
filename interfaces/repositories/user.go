package repositories

import (
	"context"

	"github.com/claudesky/identity-go/models"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	InsertUser(ctx context.Context, m *models.User) error
}
