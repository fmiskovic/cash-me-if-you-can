package api

func (r *Router) initRoutes(h handlers) {
	r.Post("/accounts", r.MakeHttpHandlerFunc(h.accountCreate.Handle))
	r.Get("/accounts/{id}", r.MakeHttpHandlerFunc(h.accountDetails.Handle))
	r.Get("/accounts", r.MakeHttpHandlerFunc(h.accountList.Handle))
	r.Post("/accounts/{id}/transactions", r.MakeHttpHandlerFunc(h.transactionCreate.Handle))
	r.Get("/accounts/{id}/transactions", r.MakeHttpHandlerFunc(h.transactionList.Handle))
	r.Post("/transfer", r.MakeHttpHandlerFunc(h.transactionTransfer.Handle))
}
