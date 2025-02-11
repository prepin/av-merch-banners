package repository

import (
	"av-merch-shop/internal/entities"
	"av-merch-shop/internal/errs"
	"av-merch-shop/pkg/database"
	"context"
	"errors"
	"fmt"
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

func NewPGUserRepo(db *database.Database) *PGUserRepo {
	return &PGUserRepo{
		db: db,
	}
}

func (u *PGUserRepo) GetByUsername(ctx context.Context, username string) (*entities.User, error) {

	stmt := psql.Select(
		sm.From("users"),
		sm.Where(psql.Raw("username").EQ(psql.Arg(username))),
		sm.Limit(1),
	)

	query, args := stmt.MustBuild(ctx)
	fmt.Println(query, args)

	row, _ := u.db.Pool.Query(ctx, query, args...)
	user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[entities.User])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound{Err: err}
		}
		u.logger.Error("Failed query user", "error", errs.ErrInternal{Err: err})
		return nil, err
	}

	return &user, nil
}

func (u *PGUserRepo) Create(ctx context.Context, data entities.UserData) (*entities.User, error) {
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
	fmt.Println(query, args)
	row, _ := u.db.Pool.Query(ctx, query, args...)
	user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[entities.User])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound{Err: err}
		}
		u.logger.Error("Failed query user", "error", errs.ErrInternal{Err: err})
		return nil, err
	}

	return &user, nil
}
