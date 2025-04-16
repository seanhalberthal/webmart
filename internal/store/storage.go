package store

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"time"
)

var (
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Products interface {
		ProductCreate(context.Context, *Product) error
		ProductGetByID(context.Context, uuid.UUID) (*Product, error)
		ProductGetAll(context.Context) ([]Product, error)
		ProductDelete(context.Context, uuid.UUID) error
		ProductUpdate(context.Context, *Product) error
	}

	Users interface {
		UserCreate(context.Context, *User) error
		UserGet(context.Context, uuid.UUID) (*User, error)
	}

	Reviews interface {
		ReviewCreate(context.Context, *Review) error
		ReviewGet(context.Context, uuid.UUID) ([]Review, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Products: &ProductStore{db},
		Users:    &UserStore{db},
		Reviews:  &ReviewStore{db},
	}
}

//func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
//	tx, err := db.BeginTx(ctx, nil)
//	if err != nil {
//		return err
//	}
//
//	if err = fn(tx); err != nil {
//		rbErr := tx.Rollback()
//		if rbErr != nil {
//			return rbErr
//		}
//		return err
//	}
//
//	return tx.Commit()
//}
