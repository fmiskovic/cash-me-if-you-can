package account

type Account struct {
	ID      string
	Owner   string
	Balance float64
}

type Option func(*Account)

func New(opts ...Option) *Account {
	a := &Account{}
	for _, opt := range opts {
		opt(a)
	}

	return a
}

func WithId(id string) Option {
	return func(a *Account) {
		a.ID = id
	}
}

func WithOwner(owner string) Option {
	return func(a *Account) {
		a.Owner = owner
	}
}

func WithBalance(balance float64) Option {
	return func(a *Account) {
		a.Balance = balance
	}
}
