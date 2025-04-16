package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Reviews     []Review  `json:"reviews"`
}

type ProductSummary struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Rating    int       `json:"rating"`
	Price     float64   `json:"price"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductStore struct {
	db *sql.DB
}

func (s *ProductStore) ProductCreate(ctx context.Context, product *Product) error {
	if product.Version == 0 {
		product.Version = 1
	}

	query := `INSERT INTO products (user_id, title, description, rating, price, stock) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	row := s.db.QueryRowContext(ctx, query, product.UserID, product.Title, product.Description, product.Rating, product.Price, product.Stock)

	err := row.Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProductStore) ProductGetByID(ctx context.Context, productID uuid.UUID) (*Product, error) {
	query := `SELECT id, user_id, title, description, price, stock, version, created_at, updated_at FROM products WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, productID)
	product := &Product{}

	err := row.Scan(
		&product.ID,
		&product.UserID,
		&product.Title,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.Version,
		&product.CreatedAt,
		&product.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}

	return product, nil
}

func (s *ProductStore) ProductGetAll(ctx context.Context) ([]ProductSummary, error) {
	query := `SELECT id, user_id, title, price, rating, created_at, updated_at FROM products`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var products []ProductSummary
	for rows.Next() {
		p := ProductSummary{}
		if err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Price,
			&p.Rating,
			&p.CreatedAt,
			&p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}

	return products, nil
}

func (s *ProductStore) ProductDelete(ctx context.Context, productID uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, productID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (s *ProductStore) ProductUpdate(ctx context.Context, product *Product) error {
	query := `UPDATE products SET title = $1, description = $2, version = version + 1 WHERE id = $3 AND version = $4 RETURNING version`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, product.Title, product.Description, product.ID, product.Version).Scan(&product.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return fmt.Errorf("product not found")
		default:
			return err
		}
	}

	return nil
}
