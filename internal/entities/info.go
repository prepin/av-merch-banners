//nolint:tagliatelle // такой формат полей в ТЗ
package entities

type UserInfo struct {
	Coins       int             `json:"coins"`
	Inventory   UserInventory   `json:"inventory"`
	CoinHistory UserCoinHistory `json:"coinHistory"`
}

type UserInventory []UserInventoryItem

type UserCoinHistory struct {
	Received UserReceived `json:"received"`
	Sent     UserSent     `json:"sent"`
}

type UserInventoryItem struct {
	Type     string `db:"codename" json:"type"`
	Quantity int    `json:"quantity"`
}

type UserSent []UserSentItem
type UserSentItem struct {
	ToUser string `db:"username" json:"toUser"`
	Amount int    `json:"amount"`
}

type UserReceived []UserReceivedItem
type UserReceivedItem struct {
	FromUser string `db:"username" json:"fromUser"`
	Amount   int    `json:"amount"`
}
