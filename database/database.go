package database

import (
	"context"
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
	rows, err = d.pool.Query(ctx, sql, args)
	return
}

func (d *PostgresDatabase) QueryRow(
	ctx context.Context,
	sql string,
	args any,
) (
	row pgx.Row,
) {
	row = d.pool.QueryRow(ctx, sql, args)
	return
}

func (d *PostgresDatabase) Exec(
	ctx context.Context,
	sql string,
	args any,
) error {
	_, err := d.pool.Exec(ctx, sql, args)
	return err
}
