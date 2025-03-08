package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type master struct {
	pool *pgxpool.Pool
}

func (db *DB) RW() *master { //nolint:revive
	return db.master
}

func (m *master) Exec(ctx context.Context, sql string, arguments ...any) error {
	_, err := m.pool.Exec(ctx, sql, arguments...)
	return err
}

func (m *master) Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error) {
	return m.pool.Query(ctx, sql, arguments...)
}

func (m *master) QueryRow(ctx context.Context, sql string, arguments ...any) pgx.Row {
	return m.pool.QueryRow(ctx, sql, arguments...)
}

// TODO: Implement transaction
func (m *master) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return m.pool.BeginTx(ctx, txOptions)
}
