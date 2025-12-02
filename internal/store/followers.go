package store

import (
	"context"
	"database/sql"
	"time"
)

type FollowStore struct {
	db *sql.DB
}

type Follower struct {
	UserID     int64     `json:"user_id"`
	FollowerID int64     `json:"follower_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (s *FollowStore) Follow(ctx context.Context, followerID int64, userID int64) error {
	query := `
		INSERT INTO followers (user_id, follower_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
	`
	_, err := s.db.ExecContext(ctx, query, userID, followerID, time.Now().UTC())
	return err
}

func (s *FollowStore) UnFollow(ctx context.Context, followerID int64, userID int64) error {
	query := `
		DELETE FROM followers
		WHERE user_id = $1 AND follower_id = $2
	`
	_, err := s.db.ExecContext(ctx, query, userID, followerID)
	return err
}
