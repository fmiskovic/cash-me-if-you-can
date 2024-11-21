package account_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/fmiskovic/cash-me-if-you-can/internal"
	"github.com/fmiskovic/cash-me-if-you-can/internal/account"
	"github.com/fmiskovic/cash-me-if-you-can/internal/account/mock"
)

func TestCreateAccount(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)

	details := &account.Details{
		AccountId: "1",
		Owner:     "Alice",
		Balance:   100.37,
	}

	tests := []struct {
		name    string
		req     account.CreateRequest
		mockFn  func(m *mock.MockRepository)
		want    *account.Details
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "create account",
			req: account.CreateRequest{
				Owner:   "Alice",
				Balance: 100.37,
			},
			mockFn: func(m *mock.MockRepository) {
				input := account.New(
					account.WithOwner("Alice"),
					account.WithBalance(100.37),
				)
				got := account.New(
					account.WithId("1"),
					account.WithOwner("Alice"),
					account.WithBalance(100.37),
				)
				m.EXPECT().Create(ctx, input).Return(got, nil)
			},
			want:    details,
			wantErr: assert.NoError,
		},
		{
			name: "create account error",
			req: account.CreateRequest{
				Owner:   "Alice",
				Balance: 100.37,
			},
			mockFn: func(m *mock.MockRepository) {
				m.EXPECT().Create(ctx, gomock.Any()).Return(nil, assert.AnError)
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}

	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mock.NewMockRepository(ctrl)

			tt.mockFn(repo)

			s := account.NewService(repo)

			got, err := s.Create(ctx, tt.req)
			if err != nil {
				tt.wantErr(t, err, "unexpected create account error")
				return
			}

			assert.Equal(t, tt.want.AccountId, got.AccountId)
			assert.Equal(t, tt.want.Owner, got.Owner)
			assert.Equal(t, tt.want.Balance, got.Balance)
		})
	}
}

func TestGetAccount(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)

	tests := []struct {
		name    string
		id      string
		mockFn  func(m *mock.MockRepository)
		want    *account.Details
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "get account",
			id:   "1",
			mockFn: func(m *mock.MockRepository) {
				got := account.New(
					account.WithId("1"),
					account.WithOwner("Alice"),
					account.WithBalance(100.37),
				)
				m.EXPECT().Get(ctx, "1").Return(got, nil)
			},
			want: &account.Details{
				AccountId: "1",
				Owner:     "Alice",
				Balance:   100.37,
			},
			wantErr: assert.NoError,
		},
		{
			name: "get account error",
			id:   "1",
			mockFn: func(m *mock.MockRepository) {
				m.EXPECT().Get(ctx, "1").Return(nil, assert.AnError)
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}

	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mock.NewMockRepository(ctrl)

			tt.mockFn(repo)

			s := account.NewService(repo)

			got, err := s.Get(ctx, tt.id)
			if err != nil {
				tt.wantErr(t, err, "unexpected get account error")
				return
			}

			assert.Equal(t, tt.want.AccountId, got.AccountId)
			assert.Equal(t, tt.want.Owner, got.Owner)
			assert.Equal(t, tt.want.Balance, got.Balance)
		})
	}
}

func TestListAccounts(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)

	items := make([]account.Details, 2)
	items[0] = account.Details{
		AccountId: "1",
		Owner:     "Alice",
		Balance:   100.37,
	}
	items[1] = account.Details{
		AccountId: "2",
		Owner:     "Bob",
		Balance:   200.37,
	}

	page := internal.Page[account.Details]{
		TotalPages: 1,
		TotalItems: 2,
		Items:      items,
	}

	tests := []struct {
		name    string
		req     internal.PageRequest
		mockFn  func(m *mock.MockRepository)
		want    internal.Page[account.Details]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "list accounts",
			req:  internal.DefaultPageRequest(),
			mockFn: func(m *mock.MockRepository) {
				acc1 := account.Account{
					ID:      "1",
					Owner:   "Alice",
					Balance: 100.37,
				}
				acc2 := account.Account{
					ID:      "2",
					Owner:   "Bob",
					Balance: 200.37,
				}

				got := internal.Page[account.Account]{
					Items:      []account.Account{acc1, acc2},
					TotalItems: 2,
					TotalPages: 1,
				}
				m.EXPECT().List(ctx, internal.DefaultPageRequest()).Return(got, nil)
			},
			want:    page,
			wantErr: assert.NoError,
		},
		{
			name: "list accounts error",
			req:  internal.DefaultPageRequest(),
			mockFn: func(m *mock.MockRepository) {
				m.EXPECT().
					List(ctx, internal.DefaultPageRequest()).
					Return(internal.Page[account.Account]{}, assert.AnError)
			},
			want:    internal.Page[account.Details]{},
			wantErr: assert.Error,
		},
	}

	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mock.NewMockRepository(ctrl)

			tt.mockFn(repo)

			s := account.NewService(repo)

			got, err := s.List(ctx, tt.req)
			if err != nil {
				tt.wantErr(t, err, "unexpected list accounts error")
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
