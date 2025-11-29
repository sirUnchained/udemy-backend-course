package store

import (
	"context"
	"database/sql"
	"time"
)

// this is user store (something like user database?)
type UserStore struct {
	db *sql.DB
}

// structure of user entity
type User struct {
	ID         int64     `json:"id"`
	UserName   string    `json:"username"`
	Email      string    `json:"email"`
	Password   string    `json:"_"`
	CreaetedAt time.Time `json:"creaeted_at"`
}

// CRUD users
func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO posts (username, email, password) VALUES($1, $2, $3, $4) RETURNING id, created_at,`

	err := s.db.QueryRowContext(ctx, query,
		user.UserName,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreaetedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
