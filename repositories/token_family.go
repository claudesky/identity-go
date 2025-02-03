package repositories

import (
	"context"

	"github.com/claudesky/identity-go/models"
	"github.com/claudesky/identity-go/services"
	"github.com/jackc/pgx/v5"
)

type TokenFamilyRepository struct {
	db *services.Database
}

func NewTokenFamilyRepository(db *services.Database) *TokenFamilyRepository {
	return &TokenFamilyRepository{db}
}

func (r *TokenFamilyRepository) GetTokensBySub(
	ctx context.Context,
	sub string,
) (
	*[]models.TokenFamily,
	error,
) {
	query := `select * from token_families where sub = @sub`
	args := pgx.NamedArgs{"sub": sub}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByNameLax[models.TokenFamily],
	)
	return &result, err
}

func (r *TokenFamilyRepository) InsertToken(
	ctx context.Context,
	m *models.TokenFamily,
) error {
	query := `insert into token_families (
		id, sub, last_issued, created_at, last_issued_at, expires_at
	) values (
		@id, @sub, @last_issued, @created_at, @last_issued_at, @expires_at
	)`
	args := pgx.NamedArgs{
		"id":             m.Id,
		"sub":            m.Sub,
		"last_issued":    m.LastIssued,
		"created_at":     m.CreatedAt,
		"last_issued_at": m.LastIssuedAt,
		"expires_at":     m.ExpiresAt,
	}

	return r.db.Exec(ctx, query, args)
}

func (r *TokenFamilyRepository) GetTokenById(
	ctx context.Context,
	id string,
) (
	*models.TokenFamily,
	error,
) {
	query := `select * from token_families where id = @id`
	args := pgx.NamedArgs{"id": id}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectOneRow(
		rows,
		pgx.RowToStructByNameLax[models.TokenFamily],
	)
	return &result, err
}

func (r *TokenFamilyRepository) UpdateToken(
	ctx context.Context,
	m *models.TokenFamily,
) error {
	query := `
		update token_families
		set
			sub = @sub,
			last_issued = @last_issued,
			created_at = @created_at,
			last_issued_at = @last_issued_at,
			expires_at = @expires_at
		where id = @id
	`

	args := pgx.NamedArgs{
		"id":             m.Id,
		"sub":            m.Sub,
		"last_issued":    m.LastIssued,
		"created_at":     m.CreatedAt,
		"last_issued_at": m.LastIssuedAt,
		"expires_at":     m.ExpiresAt,
	}

	return r.db.Exec(ctx, query, args)
}
