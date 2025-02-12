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
	RecipientID     int
	Amount          int
	TransactionType transactionType
	ReferenceId     uuid.UUID `db:"transaction_reference_id"`
	CreatedAt       time.Time
}

type TransactionData struct {
	UserID          int
	RecipientID     int
	Amount          int
	TransactionType transactionType
	ReferenceId     uuid.UUID
}
