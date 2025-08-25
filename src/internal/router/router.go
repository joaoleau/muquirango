package router

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/joaoleau/muquirango/internal/handler"
	"github.com/joaoleau/muquirango/internal/repository"
)

type Repositories struct {
	TransactionRepo *repository.TransactionRepo
}

func RegisterRoutes(r *chi.Mux, repos *Repositories) () {
	transactionHandler := handler.NewTransactionHandler(
		context.Background(),
		*repos.TransactionRepo,
	)
	
	r.Route("/api", func(r chi.Router) {
		r.Route("/transaction", func(r chi.Router) {
			r.Get("/", transactionHandler.ListTransactions)
			r.Post("/", transactionHandler.NewTransaction)
			r.Put("/{id}", transactionHandler.UpdateTransactionByID)
			r.Get("/{id}", transactionHandler.GetTransactionByID)
			r.Delete("/{id}", transactionHandler.DeleteTransactionByID)
		})
	})

}