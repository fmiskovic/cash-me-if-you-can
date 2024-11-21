package repositories

import (
	"context"
	_ "embed"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/fmiskovic/cash-me-if-you-can/database"
	"github.com/fmiskovic/cash-me-if-you-can/internal/transaction"
	"github.com/fmiskovic/cash-me-if-you-can/pkg/errorx"
)

var (
	//go:embed sql/transaction_insert.sql
	insertTransactionSql string
	//go:embed sql/account_update.sql
	updateAccountSql string
	//go:embed sql/transaction_select_by_account_id.sql
	selectTransactionsByAccountIdSql string
	//go:embed sql/transaction_select_by_id.sql
	selectTransactionByIdSql string
)

type TransactionRepository struct {
	baseRepository
}

func NewTransactionRepository(db database.Service) TransactionRepository {
	return TransactionRepository{
		baseRepository: newBaseRepository(db),
	}
}

func (r TransactionRepository) Create(ctx context.Context, t *transaction.Transaction) (*transaction.Transaction, error) {
	if t == nil {
		return nil, errorx.NewError(
			errors.New("transaction input is nil"),
			errorx.ErrInvalidInput,
		)
	}
	if t.Amount <= 0 {
		return nil, errorx.NewError(
			errors.New("invalid amount"),
			errorx.ErrInvalidInput,
		)
	}

	// execute inside transaction and rollback on error
	err := r.Execute(ctx, func(tx pgx.Tx) error {
		acc, err := r.lockAccountById(ctx, tx, t.AccountID)
		if err != nil {
			return err
		}

		// fail if account does not have enough funds
		if t.Type == transaction.Withdrawal && acc.Balance < t.Amount {
			return errorx.NewError(
				errors.New("withdrawal failed - insufficient funds"),
				errorx.ErrInvalidInput,
			)
		}

		// update account balance
		newBalance := acc.Balance + t.Amount
		if t.Type == transaction.Withdrawal {
			newBalance = acc.Balance - t.Amount
		}
		if _, err = tx.Exec(ctx, updateAccountSql, t.AccountID, newBalance); err != nil {
			return err
		}

		// create transaction
		if err = tx.QueryRow(ctx, insertTransactionSql,
			t.AccountID, // $1
			t.Amount,    // $2
			t.Type,      // $3
		).Scan(&t.ID, &t.Timestamp); err != nil {
			return err
		}

		return nil
	})

	return t, err
}

func (r TransactionRepository) Transfer(ctx context.Context, from *transaction.Transaction, to *transaction.Transaction) error {
	if from == nil || to == nil {
		return errorx.NewError(
			errors.New("transaction input is nil"),
			errorx.ErrInvalidInput,
		)
	}
	if from.Amount <= 0 || to.Amount <= 0 {
		return errorx.NewError(
			errors.New("invalid amount"),
			errorx.ErrInvalidInput,
		)
	}

	// execute inside transaction and rollback on error
	err := r.Execute(ctx, func(tx pgx.Tx) error {
		accFrom, err := r.lockAccountById(ctx, tx, from.AccountID)
		if err != nil {
			return err
		}
		accTo, err := r.lockAccountById(ctx, tx, to.AccountID)
		if err != nil {
			return err
		}

		// check if account has enough funds to make a transfer
		if accFrom.Balance < from.Amount {
			return errorx.NewError(
				errors.New("transfer failed - insufficient funds"),
				errorx.ErrInvalidInput,
			)
		}

		// update account-from balance
		accFrom.Balance -= from.Amount
		if _, err = tx.Exec(ctx, updateAccountSql, from.AccountID, accFrom.Balance); err != nil {
			return err
		}

		// update account-to balance
		accTo.Balance += to.Amount
		if _, err = tx.Exec(ctx, updateAccountSql, to.AccountID, accTo.Balance); err != nil {
			return err
		}

		// create transactions
		_, err = tx.Exec(ctx, insertTransactionSql, from.AccountID, from.Amount, from.Type)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, insertTransactionSql, to.AccountID, to.Amount, to.Type)

		return err
	})
	return err
}

func (r TransactionRepository) GetByAccountId(ctx context.Context, accountId string) ([]transaction.Transaction, error) {
	// first check if account exists
	var exist bool
	err := r.Pool().QueryRow(ctx, accountExistSql, accountId).Scan(&exist)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errorx.NewError(
			errors.New("account not found"),
			errorx.ErrNotFound,
		)
	}

	rows, err := r.Pool().Query(ctx, selectTransactionsByAccountIdSql, accountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trs []transaction.Transaction
	for rows.Next() {
		var tr transaction.Transaction
		if err = rows.Scan(&tr.ID, &tr.AccountID, &tr.Type, &tr.Amount, &tr.Timestamp); err != nil {
			return nil, err
		}
		trs = append(trs, tr)
	}

	return trs, nil
}

func (r TransactionRepository) GetById(ctx context.Context, id string) (*transaction.Transaction, error) {
	tr := new(transaction.Transaction)
	err := r.Pool().
		QueryRow(ctx, selectTransactionByIdSql, id).
		Scan(&tr.ID, &tr.AccountID, &tr.Type, &tr.Amount, &tr.Timestamp)
	return tr, err
}
