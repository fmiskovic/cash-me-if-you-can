package mappers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fmiskovic/cash-me-if-you-can/internal/account"
)

type AccountGetRequestMapper struct{}

func (m *AccountGetRequestMapper) Map(r *http.Request) (string, error) {
	id := r.PathValue("id")
	if id == "" {
		return "", errors.New("path is missing id parameter")
	}
	return id, nil
}

type AccountGetResponseMapper struct{}

func (m *AccountGetResponseMapper) Map(w http.ResponseWriter, res *account.Details) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(res)
}
