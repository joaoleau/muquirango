package dto

import (
	"github.com/joaoleau/muquirango/model"
)

type CreateEntryInput struct {
	Type        model.EntryType `json:"type"`
	Description string          `json:"description"`
}