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

	// Run migrations
	if err := goose.Up(db, absPath); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func (td *TestDatabase) CleanDatabase() error {
	db, err := sql.Open("pgx", td.ConnectionString())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Get all table names from the public schema
	rows, err := db.Query(`
        SELECT tablename
        FROM pg_tables
        WHERE schemaname = 'public'
        AND tablename != 'goose_db_version'
    `)
	if err != nil {
		return fmt.Errorf("failed to get table names: %w", err)
	}
	defer rows.Close()

	// Truncate each table
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("failed to scan table name: %w", err)
		}

		_, err = db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tableName))
		if err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", tableName, err)
		}
	}

	return rows.Err()
}
