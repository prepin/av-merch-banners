package entities

import "time"

type Order struct {
	ID            int
	UserID        int
	ItemID        int
	TransactionID int
	CreatedAt     time.Time
}

type OrderRequest struct {
	UserID   int
	ItemName string
}

type OrderData struct {
	UserID        int
	ItemID        int
	TransactionID int
}
