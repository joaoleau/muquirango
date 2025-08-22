package interfaces

import (
	"context"
	"github.com/joaoleau/muquirango/model"
)

type EntryRepository interface {
	Add (ctx context.Context, entry *model.Entry) (*model.Entry, error)
	Update (ctx context.Context, entry *model.Entry) (*model.Entry, error)
	GetAll (ctx context.Context) (*[]model.Entry, error)
	GetById (ctx context.Context, id string) (*model.Entry, error)
	Delete (ctx context.Context, entry *model.Entry) (*model.Entry, error)
}

type UserRepository interface {
	Add (ctx context.Context, entry *model.Entry) (*model.Entry, error)
	UpdateMe (ctx context.Context, entry *model.Entry) (*model.Entry, error)
	GetMe (ctx context.Context) (*model.Entry, error)
	DeleteMe (ctx context.Context, entry *model.Entry) (*model.Entry, error)
}