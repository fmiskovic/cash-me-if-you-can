package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fmiskovic/cash-me-if-you-can/internal/transaction"
)

func (s *E2ETestSuite) TestCreateTransaction() {
	tests := []struct {
		name      string
		accountId string
		input     transaction.CreateRequest
		wantCode  int
	}{
		{
			name:      "valid request",
			accountId: "c1d2e3f4-3333-4444-5555-666677778888",
			input: transaction.CreateRequest{
				Type:   transaction.Deposit,
				Amount: 100.12345,
			},
			wantCode: http.StatusCreated,
		},
		{
			name:      "invalid account id",
			accountId: "c1d2e3f4-3333-4444-5555-000000000000",
			input: transaction.CreateRequest{
				Type:   transaction.Deposit,
				Amount: 100.12345,
			},
			wantCode: http.StatusNotFound,
		},
		{
			name:      "invalid transaction type",
			accountId: "c1d2e3f4-3333-4444-5555-666677778888",
			input: transaction.CreateRequest{
				Type:   "invalid",
				Amount: 100.12345,
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name:      "invalid amount",
			accountId: "c1d2e3f4-3333-4444-5555-666677778888",
			input: transaction.CreateRequest{
				Type:   transaction.Deposit,
				Amount: -100.12345,
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name:      "missing amount",
			accountId: "c1d2e3f4-3333-4444-5555-666677778888",
			input: transaction.CreateRequest{
				Type: transaction.Deposit,
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name:      "insufficient funds",
			accountId: "c1d2e3f4-3333-4444-5555-666677778888",
			input: transaction.CreateRequest{
				Type:   transaction.Withdrawal,
				Amount: 1000000,
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name:      "missing type",
			accountId: "c1d2e3f4-3333-4444-5555-666677778888",
			input: transaction.CreateRequest{
				Amount: 100.12345,
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		tt := tc
		s.T().Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reqBody, err := json.Marshal(tt.input)
			s.NoError(err)

			path := "/accounts/" + tt.accountId + "/transactions"
			req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(reqBody))
			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			s.Equal(tt.wantCode, w.Code)

			if tt.wantCode != http.StatusCreated {
				return
			}

			var res transaction.Details
			err = json.NewDecoder(w.Body).Decode(&res)
			s.NoError(err)
			s.NotEmpty(res.TransactionId)
			s.NotEmpty(res.Timestamp)
			s.Equal(string(tt.input.Type), res.Type)
			s.Equal(tt.input.Amount, res.Amount)
			s.Equal(tt.accountId, res.AccountId)
		})
	}
}

func (s *E2ETestSuite) TestGetAccountTransactions() {
	tests := []struct {
		name      string
		accountId string
		wantCode  int
	}{
		{
			name:      "valid request",
			accountId: "b1c2d3e4-2222-3333-4444-555566667777",
			wantCode:  http.StatusOK,
		},
		{
			name:      "invalid account id",
			accountId: "b1c2d3e4-2222-3333-4444-000000000000",
			wantCode:  http.StatusNotFound,
		},
		{
			name:      "no transactions",
			accountId: "2f6f112a-a8e2-42c3-a6b0-c15e86d01704",
			wantCode:  http.StatusOK,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path := "/accounts/" + tt.accountId + "/transactions"
			req := httptest.NewRequest(http.MethodGet, path, nil)
			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			s.Equal(tt.wantCode, w.Code)

			if tt.wantCode != http.StatusOK {
				return
			}

			var res []transaction.Details
			err := json.NewDecoder(w.Body).Decode(&res)
			s.NoError(err)
		})
	}
}

func (s *E2ETestSuite) TestTransfer() {
	tests := []struct {
		name     string
		input    transaction.TransferRequest
		wantCode int
	}{
		{
			name: "valid request",
			input: transaction.TransferRequest{
				FromAccountID: "a1b2c3d4-1111-2222-3333-444455556666",
				ToAccountID:   "c1d2e3f4-3333-4444-5555-666677778888",
				Amount:        100.12345,
			},
			wantCode: http.StatusCreated,
		},
		{
			name: "invalid from account id",
			input: transaction.TransferRequest{
				FromAccountID: "a1b2c3d4-0000-2222-3333-000000000000",
				ToAccountID:   "c1d2e3f4-3333-4444-5555-666677778888",
				Amount:        100.12345,
			},
			wantCode: http.StatusNotFound,
		},
		{
			name: "invalid to account id",
			input: transaction.TransferRequest{
				FromAccountID: "a1b2c3d4-1111-2222-3333-444455556666",
				ToAccountID:   "c1d2e3f4-0000-4444-5555-000000000000",
				Amount:        100.12345,
			},
			wantCode: http.StatusNotFound,
		},
		{
			name: "insufficient funds",
			input: transaction.TransferRequest{
				FromAccountID: "a1b2c3d4-1111-2222-3333-444455556666",
				ToAccountID:   "c1d2e3f4-3333-4444-5555-666677778888",
				Amount:        1000000,
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "invalid amount",
			input: transaction.TransferRequest{
				FromAccountID: "a1b2c3d4-1111-2222-3333-444455556666",
				ToAccountID:   "c1d2e3f4-3333-4444-5555-666677778888",
				Amount:        -100.12345,
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reqBody, err := json.Marshal(tt.input)
			s.NoError(err)

			req := httptest.NewRequest(http.MethodPost, "/transfer", bytes.NewReader(reqBody))
			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			s.Equal(tt.wantCode, w.Code)

			if tt.wantCode != http.StatusCreated {
				return
			}

			var res transaction.TransferResponse
			err = json.NewDecoder(w.Body).Decode(&res)
			s.NoError(err)
			s.NotEmpty(res.ToAccountId)
			s.NotEmpty(res.FromAccountId)
			s.True(res.Amount > 0)
		})
	}
}
