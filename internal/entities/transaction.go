package entities

import (
	"time"

	"github.com/google/uuid"
)

type transactionType string

const (
	TransactionCredit   transactionType = "CREDIT"
	TransactionTransfer transactionType = "TRANSFER"
	TransactionOrder    transactionType = "ORDER"
)

type Transaction struct {
	ID              int
	UserID          int
	CounterpartyID  int
	Amount          int
	TransactionType transactionType
	ReferenceID     uuid.UUID `db:"transaction_reference_id"`
	CreatedAt       time.Time
}

type TransactionData struct {
	UserID          int
	CounterpartyID  int
	Amount          int
	TransactionType transactionType
	ReferenceID     uuid.UUID
}
