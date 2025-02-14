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
	"github.com/stephenafamo/bob/dialect/psql/sm"
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
			psql.Arg(data.ItemID),
			psql.Arg(data.TransactionID),
		),
		im.Returning("id", "user_id", "item_id", "transaction_id", "created_at"),
	)

	query, args := stmt.MustBuild(ctx)

	row, err := r.db.Conn(ctx).Query(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed create user", "error", errs.InternalError{Err: err})
		return nil, err
	}

	order, err := pgx.CollectOneRow(row, pgx.RowToStructByName[entities.Order])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFoundError{Err: err}
		}
		r.logger.Error("Failed create user", "error", errs.InternalError{Err: err})
		return nil, err
	}

	return &order, nil
}

func (r *PGOrderRepo) GetUserInventory(ctx context.Context, userID int) (*entities.UserInventory, error) {
	stmt := psql.Select(
		sm.Columns("i.codename", psql.Raw("count(o.id) as quantity")),
		sm.From("orders").As("o"),
		sm.InnerJoin("items").As("i").On(psql.Raw("o.item_id=i.id")),
		sm.Where(psql.Quote("user_id").EQ(psql.Arg(userID))),
		sm.GroupBy("item_id, i.codename"),
	)
	query, args := stmt.MustBuild(ctx)

	rows, err := r.db.Conn(ctx).Query(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed query user", "error", errs.InternalError{Err: err})
		return nil, err
	}

	inventory, err := pgx.CollectRows(rows, pgx.RowToStructByName[entities.UserInventoryItem])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFoundError{Err: err}
		}
		r.logger.Error("Failed query user", "error", errs.InternalError{Err: err})
		return nil, err
	}

	result := entities.UserInventory(inventory)
	return &result, nil
}
