package repositories

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/fmiskovic/cash-me-if-you-can/database"
	"github.com/fmiskovic/cash-me-if-you-can/internal"
	"github.com/fmiskovic/cash-me-if-you-can/internal/account"
	"github.com/fmiskovic/cash-me-if-you-can/pkg/errorx"
)

var (
	//go:embed sql/account_insert.sql
	insertAccountSql string
	//go:embed sql/account_total_count.sql
	totalAccountCountSql string
	//go:embed sql/account_select_page.sql
	accountPageSql string
	//go:embed sql/account_delete.sql
	deleteAccountSql string
)

type AccountRepository struct {
	baseRepository
}

func NewAccountRepository(db database.Service) AccountRepository {
	return AccountRepository{
		newBaseRepository(db),
	}
}

func (r AccountRepository) Get(ctx context.Context, id string) (*account.Account, error) {
	var a *account.Account
	var err error
	err = r.Execute(ctx, func(tx pgx.Tx) error { // execute opens db transaction and locks the row
		a, err = r.lockAccountById(ctx, tx, id)
		return err
	})
	return a, err
}

func (r AccountRepository) Create(ctx context.Context, acc *account.Account) (*account.Account, error) {
	if err := r.Pool().QueryRow(ctx, insertAccountSql,
		acc.Owner,   // $1
		acc.Balance, // $2
	).Scan(&acc.ID); err != nil {
		if strings.Contains(err.Error(), "accounts_owner_check") {
			return nil, errorx.NewError(
				fmt.Errorf("account with owner %s already exists", acc.Owner),
				errorx.ErrInvalidInput,
			)
		}
		return nil, err
	}

	return acc, nil
}

func (r AccountRepository) List(ctx context.Context, pageReq internal.PageRequest) (internal.Page[account.Account], error) {
	var accounts []account.Account

	rows, err := r.Pool().Query(ctx, accountPageSql, pageReq.Limit, pageReq.Offset)
	if err != nil {
		return internal.EmptyPage[account.Account](), err
	}
	defer rows.Close()

	for rows.Next() {
		var acc account.Account
		if err = rows.Scan(&acc.ID, &acc.Owner, &acc.Balance); err != nil {
			return internal.EmptyPage[account.Account](), err
		}
		accounts = append(accounts, acc)
	}

	if err = rows.Err(); err != nil {
		return internal.EmptyPage[account.Account](), err
	}

	var count int
	_ = r.Pool().QueryRow(ctx, totalAccountCountSql).Scan(&count) // ignore error, it's not critical

	totalPages := 0
	if count != 0 && pageReq.Limit != 0 {
		totalPages = (count + pageReq.Limit - 1) / pageReq.Limit
	}

	return internal.Page[account.Account]{
		Items:      accounts,
		TotalItems: count,
		TotalPages: totalPages,
	}, nil
}

func (r AccountRepository) Delete(ctx context.Context, id string) error {
	_, err := r.Pool().Exec(ctx, deleteAccountSql, id)
	return err
}
