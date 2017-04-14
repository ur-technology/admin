package user

import "ur/blockchain"

type Wallet struct {
	Address                 string                 `json:"address"`
	AnnouncementTransaction blockchain.Transaction `json:"announcementTransaction"`
	CreatedAt               Timestamp              `json:"createdAt"`
	CurrentBalance          string                 `json:"currentBalance"`
}
