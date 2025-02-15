package repository

import (
	"av-merch-shop/internal/entities"
	"av-merch-shop/internal/errs"
	"av-merch-shop/pkg/database"
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
)

type PGTransactionRepo struct {
	db     *database.Database
	logger *slog.Logger
}

func NewPGTransactionRepo(db *database.Database, logger *slog.Logger) *PGTransactionRepo {
	return &PGTransactionRepo{
		db:     db,
		logger: logger,
	}
}

// Возвращает баланс указанного пользователя. Если пользователя не существует или у него нет транзакций, вернётся 0.
func (r *PGTransactionRepo) GetUserBalance(ctx context.Context, userID int) (int, error) {

	stmt := psql.Select(
		sm.Columns(psql.F("coalesce", psql.F("sum", "amount"), 0)),
		sm.From("transactions"),
		sm.Where(psql.Quote("user_id").EQ(psql.Arg(userID))))

	query, args := stmt.MustBuild(ctx)

	var balance int
	err := r.db.Conn(ctx).QueryRow(ctx, query, args...).Scan(&balance)

	if err != nil {
		r.logger.Error("Failed query user balance", "error", err)
		return 0, err
	}

	return balance, nil
}

// Создаёт запись о новой транзакции.
func (r *PGTransactionRepo) Create(ctx context.Context, data entities.TransactionData) (*entities.Transaction, error) {

	var counterpartyID any = data.CounterpartyID
	if data.CounterpartyID == 0 {
		counterpartyID = nil
	}

	stmt := psql.Insert(
		im.Into(
			"transactions",
			"user_id", "counterparty_id",
			"amount", "transaction_type",
			"transaction_reference_id",
		),
		im.Values(
			psql.Arg(
				data.UserID, counterpartyID,
				data.Amount, data.TransactionType,
				data.ReferenceID,
			),
		),
		im.Returning("id", "user_id", "counterparty_id", "amount", "transaction_type", "transaction_reference_id", "created_at"),
	)

	query, args := stmt.MustBuild(ctx)

	row, err := r.db.Conn(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	// тут пришлось заморочиться с конверсией NULL поля
	transaction, err := pgx.CollectOneRow(row, func(row pgx.CollectableRow) (entities.Transaction, error) {
		var t entities.Transaction
		var counterpartyID sql.NullInt64

		err := row.Scan(
			&t.ID, &t.UserID, &counterpartyID,
			&t.Amount, &t.TransactionType, &t.ReferenceID,
			&t.CreatedAt,
		)
		if err != nil {
			return t, err
		}

		if counterpartyID.Valid {
			t.CounterpartyID = int(counterpartyID.Int64)
		} else {
			t.CounterpartyID = 0
		}

		return t, nil
	})

	return &transaction, err
}

func (r *PGTransactionRepo) GetOutgoingForUser(ctx context.Context, userID int) (*entities.UserSent, error) {
	stmt := psql.Select(
		sm.Columns("u.username", psql.Raw("abs(sum(t.amount)) as amount")),
		sm.From("transactions").As("t"),
		sm.InnerJoin("users").As("u").On(psql.Raw("t.counterparty_id=u.id")),
		sm.Where(psql.Raw("t.user_id").EQ(psql.Arg(userID))),
		sm.Where(psql.Raw("t.transaction_type").EQ(psql.Arg(entities.TransactionOutTransfer))),
		sm.GroupBy("u.username"),
	)
	query, args := stmt.MustBuild(ctx)

	rows, err := r.db.Conn(ctx).Query(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed query outgoing transaction list", "error", errs.InternalError{Err: err})
		return nil, err
	}

	sent, err := pgx.CollectRows(rows, pgx.RowToStructByName[entities.UserSentItem])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFoundError{Err: err}
		}
		r.logger.Error("Failed query outgoing transaction list", "error", errs.InternalError{Err: err})
		return nil, err
	}

	result := entities.UserSent(sent)
	return &result, nil
}

func (r *PGTransactionRepo) GetIncomingForUser(ctx context.Context, userID int) (*entities.UserReceived, error) {
	stmt := psql.Select(
		sm.Columns("u.username", psql.Raw("sum(t.amount) as amount")),
		sm.From("transactions").As("t"),
		sm.InnerJoin("users").As("u").On(psql.Raw("t.counterparty_id=u.id")),
		sm.Where(psql.Raw("t.user_id").EQ(psql.Arg(userID))),
		sm.Where(psql.Raw("t.transaction_type").EQ(psql.Arg(entities.TransactionInTransfer))),
		sm.GroupBy("u.username"),
	)
	query, args := stmt.MustBuild(ctx)

	rows, err := r.db.Conn(ctx).Query(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed query incoming transaction list", "error", errs.InternalError{Err: err})
		return nil, err
	}

	received, err := pgx.CollectRows(rows, pgx.RowToStructByName[entities.UserReceivedItem])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFoundError{Err: err}
		}
		r.logger.Error("Failed query incoming transaction list", "error", errs.InternalError{Err: err})
		return nil, err
	}

	result := entities.UserReceived(received)
	return &result, nil
}
