//go:generate mockgen -source=service.go -destination=./mock/service.go -package=mock
package transaction

import (
	"context"

	"github.com/softika/slogging"

	"github.com/fmiskovic/cash-me-if-you-can/pkg/errorx"
)

type Repository interface {
	Create(ctx context.Context, t *Transaction) (*Transaction, error)
	Transfer(ctx context.Context, from *Transaction, to *Transaction) error
	GetByAccountId(ctx context.Context, accountId string) ([]Transaction, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) Create(ctx context.Context, req CreateRequest) (*Details, error) {
	input := New(
		WithAccountID(req.AccountID),
		WithType(req.Type),
		WithAmount(req.Amount),
	)

	t, err := s.repo.Create(ctx, input)
	if err != nil {
		logger := slogging.Slogger()
		logger.ErrorContext(ctx, "failed to create transaction", "error", err)
		return nil, err
	}

	return &Details{
		TransactionId: t.ID,
		AccountId:     t.AccountID,
		Amount:        t.Amount,
		Timestamp:     t.Timestamp,
		Type:          string(t.Type),
	}, nil
}

func (s Service) Transfer(ctx context.Context, req TransferRequest) (*TransferResponse, error) {
	if req.FromAccountID == req.ToAccountID {
		return nil, errorx.NewErrorMsg(
			"from and to account ids are the same",
			errorx.ErrInvalidInput,
		)
	}

	from := New(
		WithAccountID(req.FromAccountID),
		WithType(Withdrawal),
		WithAmount(req.Amount),
	)

	to := New(
		WithAccountID(req.ToAccountID),
		WithType(Deposit),
		WithAmount(req.Amount),
	)

	if err := s.repo.Transfer(ctx, from, to); err != nil {
		logger := slogging.Slogger()
		logger.ErrorContext(ctx, "failed to make a transfer", "error", err)
		return nil, err
	}

	return &TransferResponse{
		FromAccountId: from.AccountID,
		ToAccountId:   to.AccountID,
		Amount:        from.Amount,
	}, nil
}

func (s Service) GetByAccountId(ctx context.Context, accountId string) ([]Details, error) {
	trs, err := s.repo.GetByAccountId(ctx, accountId)
	if err != nil {
		logger := slogging.Slogger()
		logger.ErrorContext(ctx, "failed to get list of transactions by account id", "error", err)
		return nil, err
	}

	if len(trs) == 0 {
		return nil, nil
	}

	details := make([]Details, len(trs))
	for i, tr := range trs {
		details[i] = Details{
			TransactionId: tr.ID,
			AccountId:     tr.AccountID,
			Amount:        tr.Amount,
			Timestamp:     tr.Timestamp,
			Type:          string(tr.Type),
		}
	}

	return details, nil
}
