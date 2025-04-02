package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Review struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"post_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `json:"user"`
}

type ReviewStore struct {
	db *sql.DB
}

func (s *ReviewStore) ReviewGet(ctx context.Context, postID uuid.UUID) ([]Review, error) {
	query := `SELECT r.product_id, r.user_id, r.content, r.created_at, users.username, users.id FROM reviews r JOIN users ON users.id = r.user_id
         WHERE r.product_id = $1 ORDER BY r.created_at DESC;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("Error closing rows: ", err)
			return
		}
	}(rows)

	var reviews []Review
	for rows.Next() {
		var r Review
		r.User = User{}
		err := rows.Scan(&r.ID, &r.ProductID, &r.UserID, &r.Content, &r.CreatedAt, &r.User.Username, &r.User.ID)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, r)
	}
	return reviews, nil
}

func (s *ReviewStore) ReviewCreate(ctx context.Context, r *Review) error {
	query := `INSERT INTO reviews (product_id, user_id, content) VALUES ($1, $2, $3) RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows := s.db.QueryRowContext(ctx, query, r.ProductID, r.UserID, r.Content)
	err := rows.Scan(&r.ID, &r.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
