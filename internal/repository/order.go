package repository

import (
	"av-merch-shop/internal/entities"
	"av-merch-shop/internal/errs"
	"av-merch-shop/pkg/database"
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/im"
)

type PGOrderRepo struct {
	db     *database.Database
	logger *slog.Logger
}

func NewPGOrderRepo(db *database.Database, logger *slog.Logger) *PGOrderRepo {
	return &PGOrderRepo{
		db:     db,
		logger: logger,
	}
}

func (r *PGOrderRepo) Create(ctx context.Context, data entities.OrderData) (*entities.Order, error) {
	stmt := psql.Insert(
		im.Into("orders", "user_id", "item_id", "transaction_id"),
		im.Values(
			psql.Arg(data.UserID),
			psql.Arg(data.ItemId),
			psql.Arg(data.TransactionId),
		),
		im.Returning("id", "user_id", "item_id", "transaction_id", "created_at"),
	)

	query, args := stmt.MustBuild(ctx)

	row, _ := r.db.Conn(ctx).Query(ctx, query, args...)

	order, err := pgx.CollectOneRow(row, pgx.RowToStructByName[entities.Order])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound{Err: err}
		}
		r.logger.Error("Failed create user", "error", errs.ErrInternal{Err: err})
		return nil, err
	}

	return &order, nil
}
