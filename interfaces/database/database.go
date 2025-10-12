package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Database interface {
	Ping(ctx context.Context) error
	Close()
	Query(ctx context.Context, sql string, args any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args any) pgx.Row
	Exec(ctx context.Context, sql string, args any) error
	BeginTransaction(ctx context.Context) (context.Context, error)
	CommitTransaction(ctx context.Context) error
	RollbackTransaction(ctx context.Context) error
}
