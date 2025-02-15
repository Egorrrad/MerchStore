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

type InfoResponse struct {
	Coins       *int             `json:"coins,omitempty"`
	CoinHistory *CoinHistory     `json:"coinHistory,omitempty"`
	Inventory   *[]InventoryItem `json:"inventory,omitempty"`
}
