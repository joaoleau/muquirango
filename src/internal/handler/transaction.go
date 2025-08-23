package handler

import (
	"context"
	"net/http"
	"time"
	"github.com/google/uuid"
	"github.com/go-chi/chi/v5"
	"github.com/joaoleau/muquirango/internal/dto"
	"github.com/joaoleau/muquirango/internal/model"
	"github.com/joaoleau/muquirango/internal/repository"
)

type TransactionHandler struct {
	ctx 		context.Context
	repository repository.TransactionRepo
}

func NewTransactionHandler(ctx context.Context, repo repository.TransactionRepo) *TransactionHandler {
	return &TransactionHandler{
		repository: repo,
		ctx: ctx,
	}
}

func (e *TransactionHandler) NewTransaction(w http.ResponseWriter, r *http.Request) {
	input, err := Deserialize[dto.CreateTransactionInput](r)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, err)
		return
	}

	transaction := &model.Transaction{
		ID:          uuid.NewString(),
		Type:        input.Type,
		Description: input.Description,
		CreatedAt:   time.Now().UTC(),
	}

	savedTransaction, err := e.repository.NewTransaction(e.ctx, transaction)
	if err != nil {
		ResponseWithError(w, http.StatusInternalServerError, err)
		return
	}

	ResponseWithData(w, http.StatusCreated, savedTransaction)
}

func (e *TransactionHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := e.repository.ListTransactions(e.ctx)
	if err != nil {
		ResponseWithError(w, http.StatusInternalServerError, err)
		return 
	}

	ResponseWithData(w, http.StatusAccepted, transactions) 
}

func (e *TransactionHandler) UpdateTransactionByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	updateTransaction, err := e.repository.GetTransactionByID(e.ctx, id)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, err)
		return
	}

	input, err := Deserialize[dto.CreateTransactionInput](r)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, err)
		return
	}

	newTransaction := &model.Transaction{
		ID:          updateTransaction.ID,
		Type:        input.Type,
		Description: input.Description,
		CreatedAt:   updateTransaction.CreatedAt,
	}

	_, err = e.repository.UpdateTransaction(e.ctx, newTransaction)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, err)
		return
	}

	ResponseWithData(w, http.StatusAccepted, newTransaction)
}

func (e *TransactionHandler) DeleteTransactionByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	deleteTransaction, err := e.repository.GetTransactionByID(e.ctx, id)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, err)
	}

	_, err = e.repository.DeleteTransaction(e.ctx, deleteTransaction)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, err)
		return
	}

	ResponseWithData(w, http.StatusAccepted, deleteTransaction)
}

func (e *TransactionHandler) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	Transaction, err := e.repository.GetTransactionByID(e.ctx, id)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, err)
		return
	}

	_, err = e.repository.GetTransactionByID(e.ctx, id)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, err)
		return
	}
	ResponseWithData(w, http.StatusAccepted, Transaction)
}
