package postgres

import (
	"context"

	"github.com/claudesky/identity-go/interfaces/database"
	"github.com/claudesky/identity-go/models"
	"github.com/claudesky/identity-go/utils"
	"github.com/jackc/pgx/v5"
)

type ClientRepository struct {
	db database.Database
}

func NewClientRepository(db database.Database) *ClientRepository {
	return &ClientRepository{db}
}

func (r *ClientRepository) GetClientById(
	ctx context.Context,
	id string,
) (
	*models.Client,
	error,
) {
	if err := utils.ValidateUUID(id); err != nil {
		return nil, err
	}

	query := `select * from clients where id = @id`
	args := pgx.NamedArgs{"id": id}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[models.Client])
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *ClientRepository) GetClientByClientId(
	ctx context.Context,
	clientId string,
) (
	*models.Client,
	error,
) {
	if err := utils.ValidateUUID(clientId); err != nil {
		return nil, err
	}

	query := `select * from clients where client_id = @client_id`
	args := pgx.NamedArgs{"client_id": clientId}

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[models.Client])
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *ClientRepository) InsertClient(
	ctx context.Context,
	m *models.Client,
) error {
	query := `insert into clients (
		id, client_id, client_secret, name, redirect_uris, created_at, updated_at
	) values (
		@id, @client_id, @client_secret, @name, @redirect_uris, @created_at, @updated_at
	)`
	args := pgx.NamedArgs{
		"id":            m.Id,
		"client_id":     m.ClientId,
		"client_secret": m.ClientSecret,
		"name":          m.Name,
		"redirect_uris": m.RedirectUris,
		"created_at":    m.CreatedAt,
		"updated_at":    m.UpdatedAt,
	}

	return r.db.Exec(ctx, query, args)
}

func (r *ClientRepository) GetAllClients(
	ctx context.Context,
) (
	[]*models.Client,
	error,
) {
	query := `select * from clients order by created_at desc`

	rows, err := r.db.Query(ctx, query, pgx.NamedArgs{})
	if err != nil {
		return nil, err
	}

	results, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[models.Client])
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *ClientRepository) UpdateClient(
	ctx context.Context,
	m *models.Client,
) error {
	query := `update clients set
		name = @name,
		redirect_uris = @redirect_uris,
		updated_at = @updated_at
	where id = @id`
	args := pgx.NamedArgs{
		"id":            m.Id,
		"name":          m.Name,
		"redirect_uris": m.RedirectUris,
		"updated_at":    m.UpdatedAt,
	}

	return r.db.Exec(ctx, query, args)
}

func (r *ClientRepository) DeleteClient(
	ctx context.Context,
	id string,
) error {
	if err := utils.ValidateUUID(id); err != nil {
		return err
	}

	query := `delete from clients where id = @id`
	args := pgx.NamedArgs{"id": id}

	return r.db.Exec(ctx, query, args)
}
