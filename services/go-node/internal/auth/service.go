package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type Store interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUserID(ctx context.Context, userID string) (*User, error)
	CreateUser(ctx context.Context, user *User) (int64, error)
	UpdatePassword(ctx context.Context, id int64, newHash string) error
	// optionally other helper methods
}

// Service is the auth business logic
type Service struct {
	store Store
	jwt   *Manager
}

func (s *Service) JWTManager() *Manager {
	return s.jwt
}

// NewService creates auth service
func NewService(s Store, jwt *Manager) *Service {
	return &Service{store: s, jwt: jwt}
}

func (s *Service) Authenticate(ctx context.Context, email, password string) (*User, error) {
	// login is email per your config (confirmed)
	user, err := s.store.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}
	if err := ComparePassword(user.PasswordHash, password); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

// Login uses email (primary) to authenticate and returns tokens
func (s *Service) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	if email == "" || password == "" {
		return nil, errors.New("missing credentials")
	}
	email = strings.TrimSpace(email)
	user, err := s.store.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("lookup: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("invalid login")
	}
	if err := ComparePassword(user.PasswordHash, password); err != nil {
		return nil, fmt.Errorf("invalid login")
	}
	token, err := s.jwt.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}
	return &LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		User: &User{
			UserID: user.UserID,
			Role:   user.Role,
			Name:   user.Name,
			Email:  user.Email,
		},
	}, nil
}

// WhoAmI returns user by user_id from token claim (supporting middleware)
func (s *Service) WhoAmI(ctx context.Context, userID string) (*User, error) {
	if userID == "" {
		return nil, errors.New("missing user id")
	}
	return s.store.GetByUserID(ctx, userID)
}

func (s *Service) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	// 1. Validate the refresh token and extract claims.
	claims, err := s.jwt.ParseAndValidate(refreshToken)
	if err != nil {
		return "", "", err
	}

	// 2. Load the user using the UserID from the token claims.
	user, err := s.store.GetByUserID(ctx, claims.UserID)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", errors.New("user not found")
	}

	// 3. Generate a new access token for the user.
	access, err := s.jwt.GenerateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	// 4. For now we simply return the same refresh token (token rotation can be added later if desired).
	return access, refreshToken, nil
}
