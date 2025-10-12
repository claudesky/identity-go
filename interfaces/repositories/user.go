package repositories

import (
	"context"

	"github.com/claudesky/identity-go/models"
)

type UserRepository interface {
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	InsertUser(ctx context.Context, m *models.User) error
}
