package store

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Storage struct {
	Products interface {
		ProductCreate(context.Context, *Product) error
		ProductGet(context.Context, uuid.UUID) (*Product, error)
		ProductDelete(context.Context, uuid.UUID) error
		ProductUpdate(context.Context, *Product) error
	}

	Users interface {
		UserCreate(context.Context, *User) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Products: &ProductStore{db},
		Users:    &UserStore{db},
	}
}
