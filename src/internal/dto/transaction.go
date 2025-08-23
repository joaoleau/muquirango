package dto

import (
	"github.com/joaoleau/muquirango/internal/model"
)

type CreateTransactionInput struct {
	Type        model.TransactionType `json:"type"`
	Description string          	`json:"description"`
}