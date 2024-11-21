package account

type CreateRequest struct {
	Owner   string  `json:"owner" validate:"required,min=2,max=72"`
	Balance float64 `json:"initial_balance" validate:"required,gt=0"`
}
