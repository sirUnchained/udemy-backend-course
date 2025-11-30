package store

import (
	"context"
	"database/sql"
	"time"
)

// comment database
type CommentStore struct {
	db *sql.DB
}

// comment model
type Comment struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	PostID    int64     `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// CRUD
func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `INSERT INTO comments (user_id, post_id, content) VALUES ($1, $2, $3) RETURNING id, created_at`

	err := s.db.QueryRowContext(ctx, query,
		comment.UserID,
		comment.PostID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
