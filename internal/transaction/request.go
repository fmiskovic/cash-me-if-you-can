package transaction

type CreateRequest struct {
	AccountID string  `json:"account_id" validate:"required"`
	Type      Type    `json:"type" validate:"required,oneof=deposit withdrawal"`
	Amount    float64 `json:"amount" validate:"required"`
}

type TransferRequest struct {
	FromAccountID string  `json:"from_account_id" validate:"required"`
	ToAccountID   string  `json:"to_account_id" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
}
