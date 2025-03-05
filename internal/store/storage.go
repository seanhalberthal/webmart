package store

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage struct {
	Products interface {
		Create(ctx context.Context, product *Product) error
	}

	Users interface {
		Create(ctx context.Context, user *User) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Products: &ProductStore{db},
		Users:    &UserStore{db},
	}
}
