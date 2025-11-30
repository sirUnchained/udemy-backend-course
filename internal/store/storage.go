package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrorNoRow = errors.New("record not found")
)

// this is our sorage struct which is going to save entities interfaces
type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, int64) (*Post, error)
	}
	Users interface {
		Create(context.Context, *User) error
	}
	Comments interface {
		Create(context.Context, *Comment) error
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db: db},
		Users:    &UserStore{db: db},
		Comments: &CommentStore{db: db},
	}
}
