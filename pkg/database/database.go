package database

import (
	"av-merch-shop/config"
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func NewDatabase(cfg config.DBConfig) *Database {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(cfg.GetConnectionString())
	if err != nil {
		log.Fatal("Failed to parse connection string:", err)
		return nil
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 5 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal("Failed to initialize pgx pool:", err)
		return nil
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Failed to initialize database:", err)
		return nil
	}

	return &Database{
		Pool: pool,
	}
}

func (d *Database) Close() {
	d.Pool.Close()
}
