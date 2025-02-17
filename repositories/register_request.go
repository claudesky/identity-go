package repositories

import (
	"context"

	"github.com/claudesky/identity-go/models"
	"github.com/claudesky/identity-go/services"
	"github.com/jackc/pgx/v5"
)

type RegisterRequestRepository struct {
	db *services.Database
}

func NewRegisterRequestRepository(
	db *services.Database,
) *RegisterRequestRepository {
	return &RegisterRequestRepository{db}
}

func (r *RegisterRequestRepository) InsertRegisterRequest(
	ctx context.Context,
	m *models.RegisterRequest,
) error {
	query := `insert into register_requests (
		id, email, password, created_at
	) values (
		@id, @email, @password, @created_at
	)`
	args := pgx.NamedArgs{
		"id":         m.Id,
		"email":      m.Email,
		"password":   m.Password,
		"created_at": m.CreatedAt,
	}

	return r.db.Exec(ctx, query, args)
}

func (r *RegisterRequestRepository) GetRegisterRequestById(
	ctx context.Context,
	id string,
) (
	*models.RegisterRequest,
	error,
) {
	query := `select * from register_requests where id = @id`
	args := pgx.NamedArgs{"id": id}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectOneRow(
		rows,
		pgx.RowToStructByNameLax[models.RegisterRequest],
	)
	return &result, err
}

func (r *RegisterRequestRepository) GetRegisterRequestByEmail(
	ctx context.Context,
	email string,
) (
	*models.RegisterRequest,
	error,
) {
	query := `select * from register_requests where email = @email`
	args := pgx.NamedArgs{"email": email}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectOneRow(
		rows,
		pgx.RowToStructByNameLax[models.RegisterRequest],
	)
	return &result, err
}
