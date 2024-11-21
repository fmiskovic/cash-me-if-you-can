package mappers

import (
	"encoding/json"
	"net/http"

	"github.com/fmiskovic/cash-me-if-you-can/internal/transaction"
)

type TransferRequestMapper struct{}

func (m *TransferRequestMapper) Map(r *http.Request) (transaction.TransferRequest, error) {
	var req transaction.TransferRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

type TransferResponseMapper struct{}

func (m *TransferResponseMapper) Map(w http.ResponseWriter, res *transaction.TransferResponse) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(res)
}
