package store

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"name"`
	Description string    `json:"description"`
	Rating      int       `json:"rating"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductStore struct {
	db *sql.DB
}

func (s *ProductStore) ProductCreate(ctx context.Context, product *Product) error {
	query := `INSERT INTO products (user_id, title, description, rating, price, stock) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	err := s.db.QueryRowContext(ctx, query, product.UserID, product.Title, product.Description, product.Rating, product.Price, product.Stock).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
