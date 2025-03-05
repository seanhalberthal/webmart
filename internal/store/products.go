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
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Rating      float64   `json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductStore struct {
	db *sql.DB
}

func (s *ProductStore) Create(ctx context.Context, product *Product) error {
	query := `INSERT INTO products (user_id, name, description, rating) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	err := s.db.QueryRowContext(ctx, query, product.UserID, product.Name, product.Description, product.Rating).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}
