package mappers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fmiskovic/cash-me-if-you-can/internal/transaction"
)

type TransactionListRequestMapper struct{}

func (m *TransactionListRequestMapper) Map(r *http.Request) (string, error) {
	id := r.PathValue("id")
	if id == "" {
		return "", errors.New("path is missing id parameter")
	}
	return id, nil
}

type TransactionListResponseMapper struct{}

func (m *TransactionListResponseMapper) Map(w http.ResponseWriter, res []transaction.Details) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(res)
}
