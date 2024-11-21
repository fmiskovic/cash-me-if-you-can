package mappers

import (
	"encoding/json"
	"net/http"

	"github.com/fmiskovic/cash-me-if-you-can/internal/account"
)

type AccountCreateRequestMapper struct{}

func (m *AccountCreateRequestMapper) Map(r *http.Request) (account.CreateRequest, error) {
	var req account.CreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

type AccountCreateResponseMapper struct{}

func (m *AccountCreateResponseMapper) Map(w http.ResponseWriter, res *account.Details) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(res)
}
