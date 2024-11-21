package mappers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fmiskovic/cash-me-if-you-can/internal/transaction"
)

type TransactionCreateRequestMapper struct{}

func (m *TransactionCreateRequestMapper) Map(r *http.Request) (transaction.CreateRequest, error) {
	accountId := r.PathValue("id")
	if accountId == "" {
		return transaction.CreateRequest{}, errors.New("missing id as path parameter")
	}
	var req transaction.CreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return req, err
	}
	req.AccountID = accountId
	return req, err
}

type TransactionCreateResponseMapper struct{}

func (m *TransactionCreateResponseMapper) Map(w http.ResponseWriter, res *transaction.Details) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(res)
}
