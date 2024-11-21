//go:generate mockgen -source=service.go -destination=./mock/service.go -package=mock
package account

import (
	"context"

	"github.com/softika/slogging"

	"github.com/fmiskovic/cash-me-if-you-can/internal"
)

type Repository interface {
	Create(context.Context, *Account) (*Account, error)
	Get(context.Context, string) (*Account, error)
	List(context.Context, internal.PageRequest) (internal.Page[Account], error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) Create(ctx context.Context, req CreateRequest) (*Details, error) {
	input := New(
		WithOwner(req.Owner),
		WithBalance(req.Balance),
	)

	a, err := s.repo.Create(ctx, input)
	if err != nil {
		logger := slogging.Slogger()
		logger.ErrorContext(ctx, "failed to create account", "error", err)
		return nil, err
	}

	return &Details{
		AccountId: a.ID,
		Owner:     a.Owner,
		Balance:   a.Balance,
	}, err
}

func (s Service) Get(ctx context.Context, id string) (*Details, error) {
	a, err := s.repo.Get(ctx, id)
	if err != nil {
		logger := slogging.Slogger()
		logger.ErrorContext(ctx, "failed to get account by id", "id", id, "err", err)
		return nil, err
	}
	return &Details{
		AccountId: a.ID,
		Owner:     a.Owner,
		Balance:   a.Balance,
	}, nil
}

func (s Service) List(ctx context.Context, req internal.PageRequest) (internal.Page[Details], error) {
	page, err := s.repo.List(ctx, req)
	if err != nil {
		logger := slogging.Slogger()
		logger.ErrorContext(ctx, "failed to get list of accounts", "error", err)
		return internal.EmptyPage[Details](), err
	}

	detailsList := make([]Details, len(page.Items))
	for i := range page.Items {
		detailsList[i] = Details{
			AccountId: page.Items[i].ID,
			Owner:     page.Items[i].Owner,
			Balance:   page.Items[i].Balance,
		}
	}

	return internal.Page[Details]{
		TotalPages: page.TotalPages,
		TotalItems: page.TotalItems,
		Items:      detailsList,
	}, nil
}
