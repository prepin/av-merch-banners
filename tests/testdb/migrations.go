package testdb

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func (td *TestDatabase) RunMigrations() error {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	migrationsPath := filepath.Join(currentDir, "..", "..", "schema", "migrations")

	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	db, err := sql.Open("pgx", td.ConnectionString())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	goose.SetLogger(goose.NopLogger())

	if err := goose.Up(db, absPath); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func (td *TestDatabase) CleanDatabase() error {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	migrationsPath := filepath.Join(currentDir, "..", "..", "schema", "migrations")

	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	db, err := sql.Open("pgx", td.ConnectionString())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	if err := goose.Reset(db, absPath); err != nil {
		return fmt.Errorf("failed to down migrations: %w", err)
	}

	return nil
}
