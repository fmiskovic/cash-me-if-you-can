package transaction

import "time"

type TransferResponse struct {
	FromAccountId string `json:"from_account_id"`
	ToAccountId   string `json:"to_account_id"`
	Amount        float64
}

type Details struct {
	TransactionId string    `json:"transaction_id"`
	AccountId     string    `json:"account_id"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	Timestamp     time.Time `json:"timestamp"`
}
