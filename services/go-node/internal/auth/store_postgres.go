package auth

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// AuthStore provides user persistence backed by Postgres.
type PostgresStore struct {
	db *sqlx.DB
}

// NewPostgresStore returns store backed by pg
func NewPostgresStore(pg *sqlx.DB) *PostgresStore {
	return &PostgresStore{db: pg}
}

func (s *PostgresStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := s.db.GetContext(ctx, &u, sqlGetByEmail, email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *PostgresStore) GetByUserID(ctx context.Context, userID string) (*User, error) {
	var u User
	err := s.db.GetContext(ctx, &u, sqlGetByUserID, userID)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *PostgresStore) CreateUser(ctx context.Context, user *User) (int64, error) {
	var id int64
	var created, updated string
	err := s.db.QueryRowContext(ctx, sqlCreateUser,
		user.UserID, user.Role, user.Name, user.Email, user.PasswordHash).Scan(&id, &created, &updated)
	if err != nil {
		return 0, err
	}
	user.ID = id
	return id, nil
}

func (s *PostgresStore) UpdatePassword(ctx context.Context, id int64, hash string) error {
	_, err := s.db.ExecContext(ctx, sqlUpdatePassword, hash, id)
	return err
}
