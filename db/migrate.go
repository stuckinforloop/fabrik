package db

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file" // file driver
	"github.com/jackc/pgx/v5/stdlib"
)

func (db *DB) Migrate() error {
	stdlibDB := stdlib.OpenDBFromPool(db.master.pool)

	driver, err := pgx.WithInstance(stdlibDB, &pgx.Config{})
	if err != nil {
		return fmt.Errorf("failed to create pgx driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no migrations to apply")
			return nil
		}

		return fmt.Errorf("failed to migrate: %w", err)
	}

	return nil
}
