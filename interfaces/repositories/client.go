package repositories

import (
	"context"

	"github.com/claudesky/identity-go/models"
)

type ClientRepository interface {
	GetClientById(ctx context.Context, id string) (*models.Client, error)
	GetAllClients(ctx context.Context) ([]*models.Client, error)
	InsertClient(ctx context.Context, m *models.Client) error
	UpdateClient(ctx context.Context, m *models.Client) error
	DeleteClient(ctx context.Context, id string) error
}
