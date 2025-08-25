package handler

import (
	"context"
	"net/http"
	"time"
	"go.uber.org/zap"
	"github.com/google/uuid"
	"github.com/go-chi/chi/v5"
	"github.com/joaoleau/muquirango/internal/dto"
	"github.com/joaoleau/muquirango/internal/model"
	"github.com/joaoleau/muquirango/internal/repository"
	"github.com/joaoleau/muquirango/internal/config/logger"
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
	logger.Info("Received request to create a new transaction")

	input, err := Deserialize[dto.CreateTransactionInput](r)
	if err != nil {
		logger.Error("Failed to deserialize request body", err)
		ResponseWithError(w, http.StatusBadRequest, err)
		return
	}

	transaction := &model.Transaction{
		ID:          uuid.NewString(),
		Title: input.Title,
		Type:        input.Type,
		Description: input.Description,
		Amount: input.Amount,
		CreatedAt:   time.Now().UTC(),
	}

	transaction.SetKeys()

	savedTransaction, err := e.repository.NewTransaction(e.ctx, transaction)
	if err != nil {
		logger.Error("Failed to save new transaction", err, zap.String("transaction_id", transaction.ID))
		ResponseWithError(w, http.StatusInternalServerError, err)
		return
	}

	logger.Info("Transaction created successfully", zap.String("transaction_id", savedTransaction.ID))
	ResponseWithData(w, http.StatusCreated, savedTransaction)
}

func (e *TransactionHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received request to list all transactions")
    query := r.URL.Query()
    
    startDate := query.Get("startDate"); if startDate == "" {
		startDate = time.Now().AddDate(0, 0, -3).Format("2006-01-02")
	}
	endDate := query.Get("endDate"); if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	transactions, err := e.repository.ListTransactions(e.ctx, startDate, endDate)
	if err != nil {
		logger.Error("Failed to fetch transactions", err)
		ResponseWithError(w, http.StatusInternalServerError, err)
		return 
	}

	logger.Info("Transactions retrieved successfully", zap.Int("count", len(*transactions)))
	ResponseWithData(w, http.StatusAccepted, transactions) 
}

func (e *TransactionHandler) UpdateTransactionByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	logger.Info("Received request to update transaction", zap.String("transaction_id", id))

    query := r.URL.Query()
    
	createdAt := query.Get("createdAt"); if createdAt == "" {
		createdAt = time.Now().Format("2006-01-02")
	}

	updateTransaction, err := e.repository.GetTransactionByID(e.ctx, id, createdAt)
	if err != nil {
		logger.Error("Transaction not found", err, zap.String("transaction_id", id))
		ResponseWithError(w, http.StatusBadRequest, err)
		return
	}

	input, err := Deserialize[dto.CreateTransactionInput](r)
	if err != nil {
		logger.Error("Failed to deserialize request body", err)
		ResponseWithError(w, http.StatusBadRequest, err)
		return
	}

	newTransaction := &model.Transaction{
		ID:          updateTransaction.ID,
		Title: 		 input.Title,
		Type:        input.Type,
		Description: input.Description,
		Amount:		 input.Amount,
		CreatedAt:   updateTransaction.CreatedAt,
		UpdatedAt: 	 time.Now().UTC(),
	}
	newTransaction.SetKeys()

	_, err = e.repository.UpdateTransaction(e.ctx, newTransaction)
	if err != nil {
		logger.Error("Failed to update transaction", err, zap.String("transaction_id", id))
		ResponseWithError(w, http.StatusBadRequest, err)
		return
	}

	logger.Info("Transaction updated successfully", zap.String("transaction_id", id))
	ResponseWithData(w, http.StatusAccepted, newTransaction)
}

func (e *TransactionHandler) DeleteTransactionByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	logger.Info("Received request to delete transaction", zap.String("transaction_id", id))

    query := r.URL.Query()
    
	createdAt := query.Get("createdAt"); if createdAt == "" {
		createdAt = time.Now().Format("2006-01-02")
	}

	deleteTransaction, err := e.repository.GetTransactionByID(e.ctx, id, createdAt)
	if err != nil {
		logger.Error("Transaction not found", err, zap.String("transaction_id", id))
		ResponseWithError(w, http.StatusBadRequest, err)
		return
	}

	_, err = e.repository.DeleteTransaction(e.ctx, deleteTransaction)
	if err != nil {
		logger.Error("Failed to delete transaction", err, zap.String("transaction_id", id))
		ResponseWithError(w, http.StatusBadRequest, err)
		return
	}

	logger.Info("Transaction deleted successfully", zap.String("transaction_id", id))
	ResponseWithData(w, http.StatusAccepted, deleteTransaction)
}

func (e *TransactionHandler) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	logger.Info("Received request to fetch transaction by ID", zap.String("transaction_id", id))

    query := r.URL.Query()
    
	createdAt := query.Get("createdAt"); if createdAt == "" {
		createdAt = time.Now().Format("2006-01-02")
	}

	transaction, err := e.repository.GetTransactionByID(e.ctx, id, createdAt)
	if err != nil {
		logger.Error("Transaction not found", err, zap.String("transaction_id", id))
		ResponseWithError(w, http.StatusBadRequest, err)
		return
	}

	logger.Info("Transaction retrieved successfully", zap.String("transaction_id", id))
	ResponseWithData(w, http.StatusAccepted, transaction)
}
