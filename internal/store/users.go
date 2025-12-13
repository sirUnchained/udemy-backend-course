package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// this is user store (something like user database?)
type UserStore struct {
	db *sql.DB
}

// structure of user entity
type User struct {
	ID        int64     `json:"id"`
	UserName  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"_"`
	CreatedAt time.Time `json:"creaeted_at"`
	IsActive  bool      `json:"is_active"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(pass string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &pass
	p.hash = hash

	return nil
}

func (p *password) Compare(pass string) error {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(pass))
	if err != nil {
		return err
	}

	return nil
}

// CRUD users
func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `INSERT INTO users (username, email, password) VALUES($1, $2, $3) RETURNING id, created_at;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query,
		user.UserName,
		user.Email,
		user.Password.hash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint \"users_email_key\""):
			return ErrDuplicatedEmail
		case strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint \"users_username_key\""):
			return ErrDuplicatedUsername
		default:
			return err
		}
	}

	return nil
}

func (s *UserStore) GetById(ctx context.Context, id int64) (*User, error) {
	query := `SELECT id, username, email, created_at FROM users WHERE id = $1`
	user := &User{}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.UserName,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UserStore) Update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		UPDATE posts
			SET username = $1, email = $2, is_active = $3
		WHERE id = $4
	`
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.UserName, user.Email, user.IsActive, user.ID)

	return err
}

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return withTeransaction(s.db, ctx, func(tx *sql.Tx) error {
		// cerate user
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		// create user invite
		if err := s.createUserInvitation(ctx, tx, token, invitationExp, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) Activate(ctx context.Context, token string) error {
	return withTeransaction(s.db, ctx, func(tx *sql.Tx) error {
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		user.IsActive = true

		if err := s.Update(ctx, tx, user); err != nil {
			return err
		}

		if err := s.deleteUserInvitation(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, username, email, password, created_at FROM users
		WHERE email = $1 AND is_active = true
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.UserName,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, invitationExp time.Duration, userID int64) error {
	query := `INSERT INTO user_invitations (token, user_id, expiry) VALUES ($1, $2, $3);`

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if _, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(invitationExp)); err != nil {
		return err
	}

	return nil
}

func (s *UserStore) deleteUserInvitation(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE FROM user_invitations WHERE id = $1;`

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if _, err := tx.ExecContext(ctx, query, userID); err != nil {
		return err
	}

	return nil
}

func (s *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
	SELECT u.id, u.email, u.created_at, u.is_active
	FROM users u
	JOIN users_invitations ui ON u.id = uiuser_id
	WHERE ui.token = $1 AND ui.expiry = $2
	`

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	hashed := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hashed[:])

	user := &User{}
	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&user.ID,
		&user.UserName,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return user, err
		}
	}

	return user, nil
}
