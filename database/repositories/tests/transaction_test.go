package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fmiskovic/cash-me-if-you-can/database/repositories"
	"github.com/fmiskovic/cash-me-if-you-can/internal/transaction"
)

func (s *RepositoriesTestSuite) TestCreateTransaction() {
	trRepo := repositories.NewTransactionRepository(s.dbService)
	accRepo := repositories.NewAccountRepository(s.dbService)

	tests := []struct {
		name    string
		input   *transaction.Transaction
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "valid deposit transaction",
			input: &transaction.Transaction{
				AccountID: "a1b2c3d4-1111-2222-3333-444455556666",
				Amount:    100,
				Type:      transaction.Deposit,
			},
			wantErr: assert.NoError,
		},
		{
			name: "valid withdrawal transaction",
			input: &transaction.Transaction{
				AccountID: "a1b2c3d4-1111-2222-3333-444455556666",
				Amount:    100,
				Type:      transaction.Withdrawal,
			},
			wantErr: assert.NoError,
		},
		{
			name: "insufficient balance",
			input: &transaction.Transaction{
				AccountID: "a1b2c3d4-1111-2222-3333-444455556666",
				Amount:    100000,
				Type:      transaction.Withdrawal,
			},
			wantErr: assert.Error,
		},
		{
			name: "non-existing account",
			input: &transaction.Transaction{
				AccountID: "a1b2c3d4-1111-2222-3333-000000000000",
				Amount:    100,
				Type:      transaction.Deposit,
			},
			wantErr: assert.Error,
		},
		{
			name: "empty account id",
			input: &transaction.Transaction{
				AccountID: "",
				Amount:    100,
				Type:      transaction.Deposit,
			},
			wantErr: assert.Error,
		},
		{
			name: "invalid amount",
			input: &transaction.Transaction{
				AccountID: "a1b2c3d4-1111-2222-3333-444455556666",
				Amount:    0,
				Type:      transaction.Deposit,
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			balanceBefore := 0.0
			accBefore, err := accRepo.Get(s.dbContainer.Ctx, tt.input.AccountID)
			if err == nil {
				balanceBefore = accBefore.Balance
			}

			tr, err := trRepo.Create(s.dbContainer.Ctx, tt.input)

			tt.wantErr(t, err, "Create() error = %v, wantErr %v", err, tt.wantErr)
			if err != nil {
				return
			}

			// assert transaction was created
			got, err := trRepo.GetById(s.dbContainer.Ctx, tr.ID)
			s.Assert().NoError(err)
			s.Assert().Equal(tr, got)

			// assert account balance was updated
			accAfter, err := accRepo.Get(s.dbContainer.Ctx, tt.input.AccountID)
			balanceAfter := accAfter.Balance

			if tt.input.Type == transaction.Deposit {
				s.Assert().Equal(balanceBefore+tt.input.Amount, balanceAfter)
			} else {
				s.Assert().Equal(balanceBefore-tt.input.Amount, balanceAfter)
			}
		})
	}
}

func (s *RepositoriesTestSuite) TestTransferTransaction() {
	trRepo := repositories.NewTransactionRepository(s.dbService)
	accountRepo := repositories.NewAccountRepository(s.dbService)

	tests := []struct {
		name    string
		from    *transaction.Transaction
		to      *transaction.Transaction
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "valid transfer",
			from: &transaction.Transaction{
				AccountID: "a1b2c3d4-1111-2222-3333-444455556666",
				Amount:    100.43,
				Type:      transaction.Withdrawal,
			},
			to: &transaction.Transaction{
				AccountID: "b1c2d3e4-2222-3333-4444-555566667777",
				Amount:    100.43,
				Type:      transaction.Deposit,
			},
			wantErr: assert.NoError,
		},
		{
			name: "insufficient from-account balance",
			from: &transaction.Transaction{
				AccountID: "a1b2c3d4-1111-2222-3333-444455556666",
				Amount:    100000,
				Type:      transaction.Withdrawal,
			},
			to: &transaction.Transaction{
				AccountID: "b1c2d3e4-2222-3333-4444-555566667777",
				Amount:    100000,
				Type:      transaction.Deposit,
			},
			wantErr: assert.Error,
		},
		{
			name: "non-existing from-account",
			from: &transaction.Transaction{
				AccountID: "a1b2c3d4-1111-2222-3333-000000000000",
				Amount:    100,
				Type:      transaction.Withdrawal,
			},
			to: &transaction.Transaction{
				AccountID: "b1c2d3e4-2222-3333-4444-555566667777",
				Amount:    100,
				Type:      transaction.Deposit,
			},
			wantErr: assert.Error,
		},
		{
			name: "non-existing to-account",
			from: &transaction.Transaction{
				AccountID: "a1b2c3d4-1111-2222-3333-444455556666",
				Amount:    100,
				Type:      transaction.Deposit,
			},
			to: &transaction.Transaction{
				AccountID: "a1b2c3d4-1111-2222-3333-000000000000",
				Amount:    100,
				Type:      transaction.Withdrawal,
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			fromAccBalanceBefore := 0.0
			fromAccBefore, err := accountRepo.Get(s.dbContainer.Ctx, tt.from.AccountID)
			if err == nil {
				fromAccBalanceBefore = fromAccBefore.Balance
			}

			toAccBalanceBefore := 0.0
			toAccBefore, err := accountRepo.Get(s.dbContainer.Ctx, tt.to.AccountID)
			if err == nil {
				toAccBalanceBefore = toAccBefore.Balance
			}

			err = trRepo.Transfer(s.dbContainer.Ctx, tt.from, tt.to)

			tt.wantErr(t, err, "Transfer() error = %v, wantErr %v", err, tt.wantErr)
			if err != nil {
				return
			}

			// assert from-account balance was updated
			fromAccAfter, err := accountRepo.Get(s.dbContainer.Ctx, tt.from.AccountID)
			fromAccBalanceAfter := fromAccAfter.Balance
			s.Assert().Equal(fromAccBalanceBefore-tt.from.Amount, fromAccBalanceAfter)

			// assert to-account balance was updated
			toAccAfter, err := accountRepo.Get(s.dbContainer.Ctx, tt.to.AccountID)
			toAccBalanceAfter := toAccAfter.Balance
			s.Assert().Equal(toAccBalanceBefore+tt.to.Amount, toAccBalanceAfter)
		})
	}
}
