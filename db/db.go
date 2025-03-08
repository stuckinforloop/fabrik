package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	master *master
	reader *replica
}

// Config holds database configuration
type Config struct {
	MasterDSN string `json:"master_dsn"`
	ReaderDSN string `json:"reader_dsn"`
	MaxConns  int32  `json:"max_conns"`
}

// New creates a new database instance with separate read and write connections
func New(ctx context.Context, cfg Config) (*DB, error) {
	if cfg.MaxConns == 0 {
		cfg.MaxConns = 10
	}

	masterConfig, err := pgxpool.ParseConfig(cfg.MasterDSN)
	if err != nil {
		return nil, fmt.Errorf("parsing master DSN: %w", err)
	}
	masterConfig.MaxConns = cfg.MaxConns

	readerConfig, err := pgxpool.ParseConfig(cfg.ReaderDSN)
	if err != nil {
		return nil, fmt.Errorf("parsing reader DSN: %w", err)
	}
	readerConfig.MaxConns = cfg.MaxConns

	masterPool, err := pgxpool.NewWithConfig(ctx, masterConfig)
	if err != nil {
		return nil, fmt.Errorf("creating master connection pool: %w", err)
	}

	readerPool, err := pgxpool.NewWithConfig(ctx, readerConfig)
	if err != nil {
		masterPool.Close()
		return nil, fmt.Errorf("creating reader connection pool: %w", err)
	}

	db := &DB{
		master: &master{masterPool},
		reader: &replica{readerPool},
	}

	return db, nil
}

func (db *DB) Close() error {
	db.master.pool.Close()
	db.reader.pool.Close()
	return nil
}
