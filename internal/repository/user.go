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

type PGUserRepo struct {
	db     *database.Database
	logger *slog.Logger
}

func NewPGUserRepo(db *database.Database, logger *slog.Logger) *PGUserRepo {
	return &PGUserRepo{
		db:     db,
		logger: logger,
	}
}

func (r *PGUserRepo) GetByUsername(ctx context.Context, username string) (*entities.User, error) {

	stmt := psql.Select(
		sm.From("users"),
		sm.Where(psql.Quote("username").EQ(psql.Arg(username))),
		sm.Limit(1),
	)

	query, args := stmt.MustBuild(ctx)

	row, err := r.db.Conn(ctx).Query(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed query user", "error", errs.InternalError{Err: err})
		return nil, err
	}

	user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[entities.User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFoundError{Err: err}
		}
		r.logger.Error("Failed query user", "error", errs.InternalError{Err: err})
		return nil, err
	}

	return &user, nil
}

func (r *PGUserRepo) GetByID(ctx context.Context, userID int) (*entities.User, error) {

	stmt := psql.Select(
		sm.From("users"),
		sm.Where(psql.Quote("id").EQ(psql.Arg(userID))),
		sm.Limit(1),
	)

	query, args := stmt.MustBuild(ctx)

	row, err := r.db.Conn(ctx).Query(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed query user", "error", errs.InternalError{Err: err})
		return nil, err
	}

	user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[entities.User])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFoundError{Err: err}
		}
		r.logger.Error("Failed query user", "error", errs.InternalError{Err: err})
		return nil, err
	}

	return &user, nil
}

func (r *PGUserRepo) Create(ctx context.Context, data entities.UserData) (*entities.User, error) {
	stmt := psql.Insert(
		im.Into("users", "username", "hashed_password", "role"),
		im.Values(
			psql.Arg(data.Username),
			psql.Arg(data.HashedPassword),
			psql.Arg(data.Role),
		),
		im.Returning("id", "username", "hashed_password, role"),
	)

	query, args := stmt.MustBuild(ctx)
	row, err := r.db.Conn(ctx).Query(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed query user", "error", errs.InternalError{Err: err})
		return nil, err
	}

	user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[entities.User])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFoundError{Err: err}
		}
		r.logger.Error("Failed query user", "error", errs.InternalError{Err: err})
		return nil, err
	}

	return &user, nil
}
