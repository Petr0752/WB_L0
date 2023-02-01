package db

import (
	"context"

	"github.com/Petr0752/L0/packages/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Database struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context) (*Database, error) {
	pgDSN := config.Get().PgDSN
	pool, err := pgxpool.New(ctx, pgDSN)
	if err != nil {
		return nil, errors.Wrap(err, "pgxpool.New")
	}

	return &Database{
		pool: pool,
	}, nil
}

func (db *Database) Close() {
	db.pool.Close()
}
