package testdb

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/pressly/goose/v3"
)

func (td *TestDatabase) LoadFixtures() error {
	db, err := sql.Open("pgx", td.ConnectionString())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	seedsPath := filepath.Join(currentDir, "..", "..", "schema", "seed")

	absPath, err := filepath.Abs(seedsPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	goose.SetLogger(goose.NopLogger())
	if err := goose.Up(db, absPath, goose.WithNoVersioning()); err != nil {
		return fmt.Errorf("failed to run seeds: %w", err)
	}

	return nil
}
