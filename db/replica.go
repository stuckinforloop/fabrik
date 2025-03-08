package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type replica struct {
	pool *pgxpool.Pool
}

func (db *DB) RO() *replica { //nolint:revive
	return db.reader
}

func (r *replica) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return r.pool.Query(ctx, sql, args...)
}

func (r *replica) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return r.pool.QueryRow(ctx, sql, args...)
}
