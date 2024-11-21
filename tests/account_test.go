package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fmiskovic/cash-me-if-you-can/internal"
	"github.com/fmiskovic/cash-me-if-you-can/internal/account"
)

func (s *E2ETestSuite) TestCreateAccount() {
	tests := []struct {
		name     string
		input    account.CreateRequest
		wantCode int
	}{
		{
			name: "valid request",
			input: account.CreateRequest{
				Balance: 100.12345,
				Owner:   "John Doe",
			},
			wantCode: http.StatusCreated,
		},
		{
			name: "missing owner",
			input: account.CreateRequest{
				Balance: 100.12345,
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "negative initial balance",
			input: account.CreateRequest{
				Balance: -100.12345,
				Owner:   "John Doe",
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

			req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(reqBody))
			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			s.Equal(tt.wantCode, w.Code)

			if tt.wantCode != http.StatusCreated {
				return
			}

			var resp account.Details
			err = json.NewDecoder(w.Body).Decode(&resp)
			s.NoError(err)

			s.NotEmpty(resp.AccountId)
			s.Equal(tt.input.Balance, resp.Balance)
		})
	}
}

func (s *E2ETestSuite) TestGetAccountDetails() {
	tests := []struct {
		name      string
		accountId string
		wantCode  int
		wantResp  account.Details
	}{
		{
			name:      "valid request",
			accountId: "2f6f112a-a8e2-42c3-a6b0-c15e86d01704",
			wantCode:  http.StatusOK,
			wantResp: account.Details{
				AccountId: "2f6f112a-a8e2-42c3-a6b0-c15e86d01704",
				Owner:     "David",
				Balance:   0.0000,
			},
		},
		{
			name:      "non-existing account",
			accountId: "2f6f112a-a8e2-42c3-a6b0-c15e86d01705",
			wantCode:  http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/accounts/"+tt.accountId, nil)
			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			s.Equal(tt.wantCode, w.Code)

			if tt.wantCode != http.StatusOK {
				return
			}

			var resp account.Details
			err := json.NewDecoder(w.Body).Decode(&resp)
			s.NoError(err)

			s.Equal(tt.wantResp, resp)
		})
	}
}

func (s *E2ETestSuite) TestListAccounts() {
	tests := []struct {
		name     string
		req      internal.PageRequest
		wantCode int
	}{
		{
			name:     "full page request",
			req:      internal.DefaultPageRequest(),
			wantCode: http.StatusOK,
		},
		{
			name: "two accounts per page",
			req: internal.PageRequest{
				Limit:  2,
				Offset: 0,
			},
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			queryParams := fmt.Sprintf("?limit=%d&offset=%d", tt.req.Limit, tt.req.Offset)
			req := httptest.NewRequest(http.MethodGet, "/accounts"+queryParams, nil)
			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			s.Equal(tt.wantCode, w.Code)

			if tt.wantCode != http.StatusOK {
				return
			}

			var page internal.Page[account.Details]
			err := json.NewDecoder(w.Body).Decode(&page)
			s.NoError(err)

			s.NotEmpty(page.Items)
			s.True(len(page.Items) <= tt.req.Limit)
			s.True(page.TotalItems > 0)
			s.True(page.TotalPages > 0)
		})
	}
}
