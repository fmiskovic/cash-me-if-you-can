package repositories

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/fmiskovic/cash-me-if-you-can/database"
	"github.com/fmiskovic/cash-me-if-you-can/internal/account"
	"github.com/fmiskovic/cash-me-if-you-can/pkg/errorx"
)

var (
	//go:embed sql/account_lock_by_id.sql
	lockAccountByIdSql string
	//go:embed sql/account_exist.sql
	accountExistSql string
)

type baseRepository struct {
	database.Service
	database.TxManager
}

func newBaseRepository(db database.Service) baseRepository {
	return baseRepository{
		Service:   db,
		TxManager: database.NewTxManager(db),
	}
}

func (r baseRepository) lockAccountById(ctx context.Context, tx pgx.Tx, id string) (*account.Account, error) {
	a := new(account.Account)

	row, err := tx.Query(ctx, lockAccountByIdSql, id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if !row.Next() {
		return nil, errorx.NewError(
			fmt.Errorf("account with id %s not found", id),
			errorx.ErrNotFound,
		)
	}

	if err = row.Scan(&a.ID, &a.Owner, &a.Balance); err != nil {
		return nil, fmt.Errorf("failed to scan account: %w", err)
	}

	// ensure no other rows exist
	if row.Next() {
		return nil, fmt.Errorf("multiple accounts found with id: %v", id)
	}

	return a, nil
}
