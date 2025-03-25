package store

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) UserCreate(ctx context.Context, user *User) error {
	query := `INSERT INTO users (name, username, email, password) VALUES ($1, $2, $3, $4)
	RETURNING id, created_at`

	err := s.db.QueryRowContext(ctx, query, user.Name, user.Username, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
