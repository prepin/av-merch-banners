package repository

import (
	"av-merch-shop/internal/entities"
	"av-merch-shop/pkg/database"
	"context"
	"database/sql"
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
func (r *PGTransactionRepo) GetUserBalance(ctx context.Context, userId int) (int, error) {

	stmt := psql.Select(
		sm.Columns(psql.F("coalesce", psql.F("sum", "amount"), 0)),
		sm.From("transactions"),
		sm.Where(psql.Quote("user_id").EQ(psql.Arg(userId))))

	query, args := stmt.MustBuild(ctx)

	var balance int
	err := r.db.Conn(ctx).QueryRow(ctx, query, args...).Scan(&balance)

	if err != nil {
		r.logger.Error("Failed query user balance", "error", err)
		return 0, err
	}

	return balance, nil
}

// Создаёт запись о новой транзакции
func (r *PGTransactionRepo) CreateTransaction(ctx context.Context, data entities.TransactionData) (*entities.Transaction, error) {

	var counterpartyId any = data.CounterpartyID
	if data.CounterpartyID == 0 {
		counterpartyId = nil
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
				data.UserID, counterpartyId,
				data.Amount, data.TransactionType,
				data.ReferenceId,
			),
		),
		im.Returning("id", "user_id", "counterparty_id", "amount", "transaction_type", "transaction_reference_id", "created_at"),
	)

	query, args := stmt.MustBuild(ctx)

	row, _ := r.db.Conn(ctx).Query(ctx, query, args...)

	// тут пришлось заморочиться с конверсией NULL поля
	transaction, err := pgx.CollectOneRow(row, func(row pgx.CollectableRow) (entities.Transaction, error) {
		var t entities.Transaction
		var counterpartyId sql.NullInt64

		err := row.Scan(
			&t.ID, &t.UserID, &counterpartyId,
			&t.Amount, &t.TransactionType, &t.ReferenceId,
			&t.CreatedAt,
		)
		if err != nil {
			return t, err
		}

		if counterpartyId.Valid {
			t.CounterpartyID = int(counterpartyId.Int64)
		} else {
			t.CounterpartyID = 0
		}

		return t, nil
	})

	return &transaction, err
}
