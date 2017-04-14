package user

import "ur/blockchain"

type Transaction struct {
	Amount          string                 `json:"amount"`
	Change          string                 `json:"change"`
	CreatedAt       Timestamp              `json:"createdAt"`
	CreatedBy       string                 `json:"createdBy"`
	Fee             string                 `json:"fee"`
	MessageText     string                 `json:"messageText"`
	MinedAt         int64                  `json:"minedAt"`
	ProfilePhotoURL string                 `json:"profilePhotoUrl"`
	Receiver        Reference              `json:"receiver"`
	Sender          Reference              `json:"sender"`
	SortKey         string                 `json:"sortKey"`
	Title           string                 `json:"title"`
	Type            string                 `json:"type"`
	UpdatedAt       Timestamp              `json:"updatedAt"`
	UrTransaction   blockchain.Transaction `json:"urTransaction"`
}

type Transactions map[string]Transaction
