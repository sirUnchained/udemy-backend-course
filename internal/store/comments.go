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
	User      User      `json:"user"`
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

	return err
}

func (s *CommentStore) GetCommentsByPostId(ctx context.Context, postID int64) ([]Comment, error) {
	query := `SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.username, users.id FROM comments AS c
				JOIN users ON users.id = c.user_id
				WHERE post_id = $1 
				ORDER BY c.created_at DESC;`

	rows, err := s.db.QueryContext(ctx, query, postID)
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var c Comment
		c.User = User{}
		err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.UserID,
			&c.Content,
			&c.CreatedAt,
			&c.User.UserName,
			&c.User.ID,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, err
}
