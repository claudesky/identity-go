package postgres

import (
	"context"

	"github.com/claudesky/identity-go/interfaces/database"
	"github.com/claudesky/identity-go/models"
	"github.com/claudesky/identity-go/utils"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db database.Database
}

func NewUserRepository(db database.Database) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) GetUserById(
	ctx context.Context,
	id string,
) (
	user *models.User,
	err error,
) {
	if err = utils.ValidateUUID(id); err != nil {
		return
	}

	query := `select * from users where id = @id`
	args := pgx.NamedArgs{"id": id}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return
	}

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[models.User])
	if err != nil {
		return
	}

	return &result, nil
}

func (r *UserRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (
	*models.User,
	error,
) {
	query := `select * from users where email = @email`
	args := pgx.NamedArgs{"email": email}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[models.User])
	if err != nil {
		return nil, err
	}

	return &result, err
}

func (r *UserRepository) InsertUser(
	ctx context.Context,
	m *models.User,
) error {
	query := `insert into users (
		id, password, email, email_verified_on
	) values (
		@id, @password, @email, @email_verified_on
	)`
	args := pgx.NamedArgs{
		"id":                m.Id,
		"password":          m.Password,
		"email":             m.Email,
		"email_verified_on": m.EmailVerifiedOn,
	}

	return r.db.Exec(ctx, query, args)
}
