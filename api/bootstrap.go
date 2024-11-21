package api

import (
	"github.com/go-playground/validator/v10"

	"github.com/fmiskovic/cash-me-if-you-can/config"
	"github.com/fmiskovic/cash-me-if-you-can/database"

	// core services
	"github.com/fmiskovic/cash-me-if-you-can/internal"
	"github.com/fmiskovic/cash-me-if-you-can/internal/account"
	"github.com/fmiskovic/cash-me-if-you-can/internal/transaction"

	// repositories
	repos "github.com/fmiskovic/cash-me-if-you-can/database/repositories"

	// http handler request-response mappers
	"github.com/fmiskovic/cash-me-if-you-can/api/mappers"
)

type repositories struct {
	account     account.Repository
	transaction transaction.Repository
}

func (r *Router) initRepositories(cfg config.DatabaseConfig) repositories {
	db := database.New(cfg)
	return repositories{
		account:     repos.NewAccountRepository(db),
		transaction: repos.NewTransactionRepository(db),
	}
}

type services struct {
	account     account.Service
	transaction transaction.Service
}

func (r *Router) initServices(repo repositories) services {
	return services{
		account:     account.NewService(repo.account),
		transaction: transaction.NewService(repo.transaction),
	}
}

type handlers struct {
	accountCreate  Handler[account.CreateRequest, *account.Details]
	accountDetails Handler[string, *account.Details]
	accountList    Handler[internal.PageRequest, internal.Page[account.Details]]

	transactionCreate   Handler[transaction.CreateRequest, *transaction.Details]
	transactionTransfer Handler[transaction.TransferRequest, *transaction.TransferResponse]
	transactionList     Handler[string, []transaction.Details]
}

func (r *Router) initHandlers(s services) handlers {
	vld := validator.New()

	accountCreateHandler := NewHandler(
		&mappers.AccountCreateRequestMapper{},
		&mappers.AccountCreateResponseMapper{},
		s.account.Create,
		vld,
	)

	accountDetailsHandler := NewHandler(
		&mappers.AccountGetRequestMapper{},
		&mappers.AccountGetResponseMapper{},
		s.account.Get,
		nil, //validation not needed for id as a string
	)

	accountListHandler := NewHandler(
		&mappers.AccountListRequestMapper{},
		&mappers.AccountListResponseMapper{},
		s.account.List,
		vld,
	)

	transactionCreateHandler := NewHandler(
		&mappers.TransactionCreateRequestMapper{},
		&mappers.TransactionCreateResponseMapper{},
		s.transaction.Create,
		vld,
	)

	transactionTransferHandler := NewHandler(
		&mappers.TransferRequestMapper{},
		&mappers.TransferResponseMapper{},
		s.transaction.Transfer,
		vld,
	)

	transactionListHandler := NewHandler(
		&mappers.TransactionListRequestMapper{},
		&mappers.TransactionListResponseMapper{},
		s.transaction.GetByAccountId,
		nil, //validation not needed for id as a string
	)

	return handlers{
		accountCreate:       accountCreateHandler,
		accountDetails:      accountDetailsHandler,
		accountList:         accountListHandler,
		transactionCreate:   transactionCreateHandler,
		transactionTransfer: transactionTransferHandler,
		transactionList:     transactionListHandler,
	}
}
