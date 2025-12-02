package store

import (
	"context"
	"database/sql"
	"errors"
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
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentsCount int `json:"comments_count"`
}

// CRUD users
func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (content, title, user_id, tags) VALUES($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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

func (s *PostStore) GetById(ctx context.Context, id int64) (*Post, error) {
	var post Post
	query := `SELECT pg_sleep(5) id, content, title, user_id, tags, version, created_at, updated_at FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.Version,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, ErrorNoRow
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, id int64) ([]PostWithMetadata, error) {
	query := `
	SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags, 
		u.username, COUNT(*) AS comments_count
	FROM posts AS p
	LEFT JOIN comments AS c on c.post_id = p.id
	LEFT JOIN users AS u ON p.user_id = u.id
	JOIN followers AS f ON f.follower_id = p.user_id OR p.user_id = $1
	WHERE f.user_id = $1 OR p.user_id = $1
	GROUP BY p.id, u.username
	ORDER BY p.created_at DESC;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []PostWithMetadata
	for rows.Next() {
		var p PostWithMetadata

		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.Version,
			pq.Array(&p.Tags),
			&p.User.UserName,
			&p.CommentsCount,
		)
		if err != nil {
			return nil, err
		}

		feed = append(feed, p)
	}

	return feed, nil
}

func (s *PostStore) DeleteById(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrorNoRow
	}

	return nil
}

func (s *PostStore) UpdateById(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts
			SET title = $1, content = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query,
		post.Title,
		post.Content,
		post.ID,
		post.Version,
	).Scan(
		&post.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrorNoRow
		default:
			return err
		}
	}

	return nil
}
