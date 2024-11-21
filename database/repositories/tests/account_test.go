package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fmiskovic/cash-me-if-you-can/database/repositories"
	"github.com/fmiskovic/cash-me-if-you-can/internal"
	"github.com/fmiskovic/cash-me-if-you-can/internal/account"
)

func (s *RepositoriesTestSuite) TestGetAccount() {
	repo := repositories.NewAccountRepository(s.dbService)

	tests := []struct {
		name    string
		id      string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "valid account",
			id:      "a1b2c3d4-1111-2222-3333-444455556666",
			wantErr: assert.NoError,
		},
		{
			name:    "non-existing account",
			id:      "a1b2c3d4-1111-2222-3333-000000000000",
			wantErr: assert.Error,
		},
		{
			name:    "empty id",
			id:      "",
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			acc, err := repo.Get(s.dbContainer.Ctx, tt.id)

			tt.wantErr(t, err, "Get() error = %v, wantErr %v", err, tt.wantErr)
			if err != nil {
				return
			}

			s.Assert().NotEmpty(acc.ID)
		})
	}
}

func (s *RepositoriesTestSuite) TestListAccounts() {
	repo := repositories.NewAccountRepository(s.dbService)

	tests := []struct {
		name    string
		input   internal.PageRequest
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "two accounts per page",
			input: internal.PageRequest{
				Limit:  2,
				Offset: 0,
			},
			wantErr: assert.NoError,
		},
		{
			name: "all accounts at the page",
			input: internal.PageRequest{
				Limit:  10,
				Offset: 0,
			},
			wantErr: assert.NoError,
		},
		{
			name: "empty page",
			input: internal.PageRequest{
				Limit:  10,
				Offset: 10,
			},
			wantErr: assert.NoError,
		},
		{
			name:    "empty page when page request is empty",
			input:   internal.PageRequest{},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			page, err := repo.List(s.dbContainer.Ctx, tt.input)

			tt.wantErr(t, err, "List() error = %v, wantErr %v", err, tt.wantErr)
			if err != nil {
				return
			}

			s.Assert().NotEmpty(page)
			s.Assert().True(len(page.Items) <= tt.input.Limit)
			s.Assert().True(page.TotalItems > 0)
		})
	}
}

func (s *RepositoriesTestSuite) TestCreateAccount() {
	repo := repositories.NewAccountRepository(s.dbService)

	tests := []struct {
		name    string
		input   *account.Account
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "valid account",
			input: &account.Account{
				Owner:   "Mark",
				Balance: 73.444416,
			},
			wantErr: assert.NoError,
		},
		{
			name:    "empty input",
			input:   &account.Account{},
			wantErr: assert.Error,
		},
		{
			name: "empty owner",
			input: &account.Account{
				Owner:   " ",
				Balance: 11.0001,
			},
			wantErr: assert.Error,
		},
		{
			name: "existing owner",
			input: &account.Account{
				Owner:   "Alice",
				Balance: 11.0001,
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			acc, err := repo.Create(s.dbContainer.Ctx, tt.input)

			tt.wantErr(t, err, "Create() error = %v, wantErr %v", err, tt.wantErr)
			if err != nil {
				return
			}

			s.Assert().NotEmpty(acc.ID)
			s.Assert().Equal(tt.input.Owner, acc.Owner)
			s.Assert().Equal(tt.input.Balance, acc.Balance)

			// assert if balance is correctly stored in the database
			got, err := repo.Get(s.dbContainer.Ctx, acc.ID)
			s.Assert().NoError(err)
			s.Assert().Equal(tt.input.Balance, got.Balance)

			// cleanup
			err = repo.Delete(s.dbContainer.Ctx, acc.ID)
			s.Assert().NoError(err)
		})
	}
}
