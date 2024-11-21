package transaction_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/fmiskovic/cash-me-if-you-can/internal/transaction"
	"github.com/fmiskovic/cash-me-if-you-can/internal/transaction/mock"
)

func TestCreateTransaction(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)

	req := transaction.CreateRequest{
		AccountID: "1",
		Type:      transaction.Deposit,
		Amount:    23.5,
	}

	tests := []struct {
		name    string
		req     transaction.CreateRequest
		mockFn  func(*mock.MockRepository)
		want    *transaction.Details
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "create transaction",
			req:  req,
			mockFn: func(repo *mock.MockRepository) {
				input := transaction.New(
					transaction.WithAccountID(req.AccountID),
					transaction.WithType(req.Type),
					transaction.WithAmount(req.Amount),
				)
				repo.EXPECT().Create(ctx, input).Return(&transaction.Transaction{
					ID:        "1",
					AccountID: "1",
					Type:      transaction.Deposit,
					Amount:    23.5,
					Timestamp: time.Now(),
				}, nil)
			},
			want: &transaction.Details{
				TransactionId: "1",
				AccountId:     "1",
				Amount:        23.5,
				Type:          string(transaction.Deposit),
			},
		},
		{
			name: "create transaction error",
			req:  req,
			mockFn: func(repo *mock.MockRepository) {
				input := transaction.New(
					transaction.WithAccountID(req.AccountID),
					transaction.WithType(req.Type),
					transaction.WithAmount(req.Amount),
				)
				repo.EXPECT().Create(ctx, input).Return(nil, assert.AnError)
			},
			wantErr: assert.Error,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mock.NewMockRepository(ctrl)

			tt.mockFn(repo)

			s := transaction.NewService(repo)

			got, err := s.Create(ctx, tt.req)
			if err != nil {
				tt.wantErr(t, err, "unexpected error")
				return
			}

			assert.Equal(t, tt.want.AccountId, got.AccountId)
			assert.Equal(t, tt.want.Amount, got.Amount)
			assert.Equal(t, tt.want.TransactionId, got.TransactionId)
			assert.Equal(t, tt.want.Type, got.Type)
		})
	}
}

func TestTransfer(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)

	req := transaction.TransferRequest{
		FromAccountID: "1",
		ToAccountID:   "2",
		Amount:        23.5,
	}

	tests := []struct {
		name    string
		req     transaction.TransferRequest
		mockFn  func(*mock.MockRepository)
		want    *transaction.TransferResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "transfer",
			req:  req,
			mockFn: func(repo *mock.MockRepository) {
				from := transaction.New(
					transaction.WithAccountID(req.FromAccountID),
					transaction.WithType(transaction.Withdrawal),
					transaction.WithAmount(req.Amount),
				)
				to := transaction.New(
					transaction.WithAccountID(req.ToAccountID),
					transaction.WithType(transaction.Deposit),
					transaction.WithAmount(req.Amount),
				)
				repo.EXPECT().Transfer(ctx, from, to).Return(nil)
			},
			want: &transaction.TransferResponse{
				FromAccountId: req.FromAccountID,
				ToAccountId:   req.ToAccountID,
				Amount:        req.Amount,
			},
			wantErr: assert.NoError,
		},
		{
			name: "transfer error",
			req:  req,
			mockFn: func(repo *mock.MockRepository) {
				from := transaction.New(
					transaction.WithAccountID(req.FromAccountID),
					transaction.WithType(transaction.Withdrawal),
					transaction.WithAmount(req.Amount),
				)
				to := transaction.New(
					transaction.WithAccountID(req.ToAccountID),
					transaction.WithType(transaction.Deposit),
					transaction.WithAmount(req.Amount),
				)
				repo.EXPECT().Transfer(ctx, from, to).Return(assert.AnError)
			},
			wantErr: assert.Error,
		},
		{
			name: "invalid transfer request",
			req: transaction.TransferRequest{
				FromAccountID: "1",
				ToAccountID:   "1",
				Amount:        23.5,
			},
			mockFn:  func(repo *mock.MockRepository) {},
			wantErr: assert.Error,
		},
	}

	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mock.NewMockRepository(ctrl)

			tt.mockFn(repo)

			s := transaction.NewService(repo)

			got, err := s.Transfer(ctx, tt.req)
			if err != nil {
				tt.wantErr(t, err, "unexpected error")
				return
			}

			assert.Equal(t, tt.want.FromAccountId, got.FromAccountId)
			assert.Equal(t, tt.want.ToAccountId, got.ToAccountId)
			assert.Equal(t, tt.want.Amount, got.Amount)
		})
	}
}

func TestGetTransactionsByAccountId(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)

	tests := []struct {
		name      string
		accountID string
		mockFn    func(*mock.MockRepository)
		want      []transaction.Details
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name:      "get transactions by account id",
			accountID: "1",
			mockFn: func(repo *mock.MockRepository) {
				repo.EXPECT().GetByAccountId(ctx, "1").Return([]transaction.Transaction{
					{
						ID:        "1",
						AccountID: "1",
						Type:      transaction.Withdrawal,
						Amount:    23.5,
						Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}, nil)
			},
			want: []transaction.Details{
				{
					TransactionId: "1",
					AccountId:     "1",
					Amount:        23.5,
					Type:          string(transaction.Withdrawal),
					Timestamp:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:      "get transactions by account id error",
			accountID: "1",
			mockFn: func(repo *mock.MockRepository) {
				repo.EXPECT().GetByAccountId(ctx, "1").Return(nil, assert.AnError)
			},
			wantErr: assert.Error,
		},
	}

	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mock.NewMockRepository(ctrl)

			tt.mockFn(repo)

			s := transaction.NewService(repo)

			got, err := s.GetByAccountId(ctx, tt.accountID)
			if err != nil {
				tt.wantErr(t, err, "unexpected error")
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
