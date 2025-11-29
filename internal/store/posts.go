package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

// this is post store (something like post database?)
type PostStore struct {
	db *sql.DB
}

// structure of post entity
type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CRUD users
func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (content, title, user_id, tags) VALUES($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	err := s.db.QueryRowContext(ctx, query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
