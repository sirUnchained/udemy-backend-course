package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrorNoRow           = errors.New("record not found")
	QueryTimeoutDuration = time.Second * 5
)

// this is our sorage struct which is going to save entities interfaces
type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, int64) (*Post, error)
		DeleteById(context.Context, int64) error
		UpdateById(context.Context, *Post) error
	}
	Users interface {
		Create(context.Context, *User) error
		GetById(context.Context, int64) (*User, error)
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetCommentsByPostId(context.Context, int64) ([]Comment, error)
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db: db},
		Users:    &UserStore{db: db},
		Comments: &CommentStore{db: db},
	}
}
