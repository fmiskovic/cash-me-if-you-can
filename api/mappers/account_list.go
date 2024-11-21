package mappers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fmiskovic/cash-me-if-you-can/internal"
	"github.com/fmiskovic/cash-me-if-you-can/internal/account"
)

type AccountListRequestMapper struct{}

func (m *AccountListRequestMapper) Map(r *http.Request) (internal.PageRequest, error) {
	offset := 0
	offsetStr := r.URL.Query().Get("offset")
	if offsetStr != "" {
		offset, _ = strconv.Atoi(offsetStr)
	}

	limit := 10
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	return internal.PageRequest{
		Limit:  limit,
		Offset: offset,
	}, nil
}

type AccountListResponseMapper struct{}

func (m *AccountListResponseMapper) Map(w http.ResponseWriter, page internal.Page[account.Details]) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(page)
}
