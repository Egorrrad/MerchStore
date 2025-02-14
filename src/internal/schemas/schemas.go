package schemas

type InventoryItem struct {
	Quantity *int    `json:"quantity,omitempty"`
	Type     *string `json:"type,omitempty"`
}

type Transaction struct {
	Amount   *int    `json:"amount,omitempty"`
	FromUser *string `json:"fromUser,omitempty"`
	ToUser   *string `json:"toUser,omitempty"`
}

type CoinHistory struct {
	Received *[]Transaction `json:"received,omitempty"`
	Sent     *[]Transaction `json:"sent,omitempty"`
}

// InfoResponse defines model for InfoResponse.
type InfoResponse struct {
	CoinHistory *CoinHistory     `json:"coinHistory,omitempty"`
	Coins       *int             `json:"coins,omitempty"`
	Inventory   *[]InventoryItem `json:"inventory,omitempty"`
}
