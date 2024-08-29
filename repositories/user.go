package repositories

import (
	"context"

	"github.com/claudesky/identity-go/models"
	"github.com/claudesky/identity-go/services"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *services.Database
}

func NewUserRepository(db *services.Database) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (
	models.User,
	error,
) {
	query := `select * from users where email = @email`
	args := pgx.NamedArgs{"email": email}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return models.User{}, err
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[models.User])
}
