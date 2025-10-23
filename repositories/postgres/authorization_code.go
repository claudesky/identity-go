package postgres

import (
	"context"

	"github.com/claudesky/identity-go/interfaces/database"
	"github.com/claudesky/identity-go/models"
	"github.com/jackc/pgx/v5"
)

type AuthorizationCodeRepository struct {
	db database.Database
}

func NewAuthorizationCodeRepository(db database.Database) *AuthorizationCodeRepository {
	return &AuthorizationCodeRepository{db}
}

func (r *AuthorizationCodeRepository) GetAuthorizationCodeByCode(
	ctx context.Context,
	code string,
) (
	*models.AuthorizationCode,
	error,
) {
	query := `select * from authorization_codes where code = @code`
	args := pgx.NamedArgs{"code": code}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[models.AuthorizationCode])
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *AuthorizationCodeRepository) InsertAuthorizationCode(
	ctx context.Context,
	m *models.AuthorizationCode,
) error {
	query := `insert into authorization_codes (
		id, code, client_id, user_id, redirect_uri, scope,
		expires_at, created_at, used, code_challenge, code_challenge_method
	) values (
		@id, @code, @client_id, @user_id, @redirect_uri, @scope,
		@expires_at, @created_at, @used, @code_challenge, @code_challenge_method
	)`
	args := pgx.NamedArgs{
		"id":                    m.Id,
		"code":                  m.Code,
		"client_id":             m.ClientId,
		"user_id":               m.UserId,
		"redirect_uri":          m.RedirectUri,
		"scope":                 m.Scope,
		"expires_at":            m.ExpiresAt,
		"created_at":            m.CreatedAt,
		"used":                  m.Used,
		"code_challenge":        m.CodeChallenge,
		"code_challenge_method": m.CodeChallengeMethod,
	}

	return r.db.Exec(ctx, query, args)
}

func (r *AuthorizationCodeRepository) MarkCodeAsUsed(
	ctx context.Context,
	code string,
) error {
	query := `update authorization_codes set used = true where code = @code`
	args := pgx.NamedArgs{"code": code}

	return r.db.Exec(ctx, query, args)
}
