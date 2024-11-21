package account

type Details struct {
	AccountId string  `json:"account_id"`
	Owner     string  `json:"owner"`
	Balance   float64 `json:"balance"`
}
