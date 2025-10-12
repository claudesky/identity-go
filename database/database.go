package database

import (
	"context"
	"log/slog"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDatabase struct {
	pool *pgxpool.Pool
}

var (
	db   *PostgresDatabase
	once sync.Once
)

func NewPostgresDatabase(
	ctx context.Context,
	cstr string,
	pass *string,
) (*PostgresDatabase, error) {
	var err error
	once.Do(func() {
		var conf *pgxpool.Config
		conf, err = pgxpool.ParseConfig(cstr)
		if err != nil {
			return
		}

		if pass != nil {
			conf.ConnConfig.Password = *pass
		}

		// // Possibly use this later? try string for now
		// conf.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		// 	pgxuuid.Register(conn.TypeMap())
		// 	return nil
		// }

		var pool *pgxpool.Pool
		pool, err = pgxpool.NewWithConfig(ctx, conf)
		if err != nil {
			return
		}

		db = &PostgresDatabase{pool}

		err = db.Ping(ctx)
	})

	return db, err
}

func (d *PostgresDatabase) Ping(ctx context.Context) error {
	return d.pool.Ping(ctx)
}

func (d *PostgresDatabase) Close() {
	d.pool.Close()
}

func (d *PostgresDatabase) Query(
	ctx context.Context,
	sql string,
	args any,
) (
	rows pgx.Rows,
	err error,
) {
	tx, ok := ctx.Value("transaction").(pgx.Tx)
	if !ok {
		rows, err = d.pool.Query(ctx, sql, args)
		return
	}
	rows, err = tx.Query(ctx, sql, args)
	return
}

func (d *PostgresDatabase) QueryRow(
	ctx context.Context,
	sql string,
	args any,
) (
	row pgx.Row,
) {
	tx, ok := ctx.Value("transaction").(pgx.Tx)
	if !ok {
		row = d.pool.QueryRow(ctx, sql, args)
		return
	}
	row = tx.QueryRow(ctx, sql, args)
	return
}

func (d *PostgresDatabase) Exec(
	ctx context.Context,
	sql string,
	args any,
) error {
	tx, ok := ctx.Value("transaction").(pgx.Tx)
	if !ok {
		_, err := d.pool.Exec(ctx, sql, args)
		return err
	} else {
		_, err := tx.Exec(ctx, sql, args)
		return err
	}
}

type contextKey string

const transactionContextKey contextKey = "transaction"

func (d *PostgresDatabase) BeginTransaction(
	ctx context.Context,
) (
	nctx context.Context,
	err error,
) {
	tx, err := d.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return
	}

	return context.WithValue(ctx, transactionContextKey, tx), nil
}

func (d *PostgresDatabase) CommitTransaction(ctx context.Context) error {
	tx, ok := ctx.Value(transactionContextKey).(pgx.Tx)
	if !ok {
		slog.Warn("Transaction Commit called with no Transaction")
		return nil // No transaction to commit
	}
	return tx.Commit(ctx)
}

func (d *PostgresDatabase) RollbackTransaction(ctx context.Context) error {
	tx, ok := ctx.Value(transactionContextKey).(pgx.Tx)
	if !ok {
		slog.Warn("Transaction Rollback called with no Transaction")
		return nil // No transaction to rollback
	}
	return tx.Rollback(ctx)
}
