package blockchain

type TransactionNonce interface{}

type Transaction struct {
	BlockHash        string           `json:"blockHash,omitempty"`
	BlockNumber      int              `json:"blockNumber"`
	From             string           `json:"from,omitempty"`
	Gas              int              `json:"gas,omitempty"`
	GasPrice         string           `json:"gasPrice,omitempty"`
	Hash             string           `json:"hash"`
	Input            string           `json:"input,omitempty"`
	Nonce            TransactionNonce `json:"nonce,omitempty"`
	R                string           `json:"r,omitempty"`
	S                string           `json:"s,omitempty"`
	To               string           `json:"to,omitempty"`
	TransactionIndex int              `json:"transactionIndex,omitempty"`
	V                string           `json:"v,omitempty"`
	Value            string           `json:"value,omitempty"`
}
