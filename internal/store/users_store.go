package store

import (
	"context"
	"database/sql"
	"errors"
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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, user.Name, user.Username, user.Email, user.Password)
	err := row.Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) UserGet(ctx context.Context, userID uuid.UUID) (*User, error) {
	query := `SELECT id, username, email, password, created_at FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	row := s.db.QueryRowContext(ctx, query, userID)

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, err
		default:
			return user, nil
		}
	}

	return user, nil
}
