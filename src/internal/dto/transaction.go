package dto

import (
	"github.com/joaoleau/muquirango/internal/model"
)

type CreateTransactionInput struct {
	Type        model.TransactionType `json:"type"`
	Title       string                `json:"title"`
	Description *string               `json:"description,omitempty"`
	Amount      int                   `json:"amount"`
}
