package database

import (
	"av-merch-shop/config"
	"context"
	"log"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	TransactionManager *manager.Manager
	pool               *pgxpool.Pool
	getter             *trmpgx.CtxGetter
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

	trManager := manager.Must(
		trmpgx.NewDefaultFactory(pool),
		manager.WithSettings(
			trmpgx.MustSettings(
				settings.Must(),
				trmpgx.WithTxOptions(pgx.TxOptions{
					IsoLevel: pgx.RepeatableRead,
				}),
			),
		),
	)

	return &Database{
		TransactionManager: trManager,
		pool:               pool,
		getter:             trmpgx.DefaultCtxGetter,
	}
}

func (d *Database) Conn(ctx context.Context) trmpgx.Tr {
	return d.getter.DefaultTrOrDB(ctx, d.pool)
}

func (d *Database) Close() {
	d.pool.Close()
}
