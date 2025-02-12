package entities

import "github.com/google/uuid"

type CreditTransactionResult struct {
	ReferenceID uuid.UUID
	NewAmount   int
}

type CreditData struct {
	Username string
	Amount   int
}
