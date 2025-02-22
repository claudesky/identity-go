package repositories

import (
	"context"

	"github.com/claudesky/identity-go/database"
	"github.com/claudesky/identity-go/models"
	"github.com/jackc/pgx/v5"
)

type EmailVerificationRepository struct {
	db *database.Database
}

func NewEmailVerificationRepository(
	db *database.Database,
) *EmailVerificationRepository {
	return &EmailVerificationRepository{db}
}

func (r *EmailVerificationRepository) InsertEmailVerificationRequest(
	ctx context.Context,
	m *models.EmailVerificationRequest,
) error {
	query := `insert into email_verification_requests (
		id, register_request_id, email, token, revoked, accepted, expires_at,
		created_at
	) values (
		@id, @register_request_id, @email, @token, @revoked, @accepted, @expires_at,
		@created_at
	)`
	args := pgx.NamedArgs{
		"id":                  m.Id,
		"register_request_id": m.RegisterRequestId,
		"email":               m.Email,
		"token":               m.Token,
		"revoked":             m.Revoked,
		"accepted":            m.Accepted,
		"expires_at":          m.ExpiresAt,
		"created_at":          m.CreatedAt,
	}

	return r.db.Exec(ctx, query, args)
}

func (
	r *EmailVerificationRepository,
) GetEmailVerificationRequestsByRegisterRequestId(
	ctx context.Context,
	id string,
) (
	*[]models.EmailVerificationRequest,
	error,
) {
	query := `select * from email_verification_requests where register_request_id = @id`
	args := pgx.NamedArgs{"id": id}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByNameLax[models.EmailVerificationRequest],
	)
	if err != nil {
		return nil, err
	}

	return &result, err
}

func (
	r *EmailVerificationRepository,
) GetEmailVerificationRequestById(
	ctx context.Context,
	id string,
) (*models.EmailVerificationRequest, error) {
	query := `select * from email_verification_requests where id = @id`
	args := pgx.NamedArgs{"id": id}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectOneRow(
		rows,
		pgx.RowToStructByNameLax[models.EmailVerificationRequest],
	)
	if err != nil {
		return nil, err
	}

	return &result, err
}

func (
	r *EmailVerificationRepository,
) Revoke(
	ctx context.Context,
	m *models.EmailVerificationRequest,
) error {
	query := `
		update email_verification_requests
		set revoked = true
		where id = @id
	`

	args := pgx.NamedArgs{
		"id": m.Id,
	}

	return r.db.Exec(ctx, query, args)
}

func (
	r *EmailVerificationRepository,
) Accept(
	ctx context.Context,
	m *models.EmailVerificationRequest,
) error {
	query := `
		update email_verification_requests
		set accepted = true
		where id = @id
	`

	args := pgx.NamedArgs{
		"id": m.Id,
	}

	return r.db.Exec(ctx, query, args)
}
