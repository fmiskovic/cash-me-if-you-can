package transaction

import "time"

type Transaction struct {
	ID        string
	AccountID string
	Type      Type
	Amount    float64
	Timestamp time.Time
}

type Type string

const (
	Deposit    Type = "deposit"
	Withdrawal Type = "withdrawal"
)

type Option func(*Transaction)

func New(opts ...Option) *Transaction {
	t := &Transaction{}
	for _, opt := range opts {
		opt(t)
	}

	return t
}

func WithAccountID(accountID string) Option {
	return func(t *Transaction) {
		t.AccountID = accountID
	}
}

func WithType(tp Type) Option {
	return func(t *Transaction) {
		t.Type = tp
	}
}

func WithAmount(amount float64) Option {
	return func(t *Transaction) {
		t.Amount = amount
	}
}
