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
	"github.com/stephenafamo/bob/dialect/psql/sm"
)

type PGItemRepo struct {
	db     *database.Database
	logger *slog.Logger
}

func NewPGItemRepo(db *database.Database, logger *slog.Logger) *PGItemRepo {
	return &PGItemRepo{
		db:     db,
		logger: logger,
	}
}

// Возвращает товар по названию
func (r *PGItemRepo) GetItemByName(ctx context.Context, itemName string) (*entities.Item, error) {

	stmt := psql.Select(
		sm.Columns("id", "codename", "cost"),
		sm.From("items"),
		sm.Where(psql.Quote("codename").EQ(psql.Arg(itemName))),
		sm.Limit(1),
	)

	query, args := stmt.MustBuild(ctx)

	row, _ := r.db.Conn(ctx).Query(ctx, query, args...)
	item, err := pgx.CollectOneRow(row, pgx.RowToStructByName[entities.Item])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound{Err: err}
		}
		r.logger.Error("Failed query item", "error", errs.ErrInternal{Err: err})
		return nil, err
	}

	return &item, nil
}
