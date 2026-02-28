package auth

import "time"

type User struct {
	ID           int64     `db:"id" json:"id"`
	UserID       string    `db:"user_id" json:"user_id"` // login username (could be email)
	Email        string    `db:"email" json:"email"`
	Role         string    `db:"role" json:"role"`
	Name         string    `db:"name" json:"name"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresInSec int64  `json:"expires_in"`
	// user info
	User struct {
		UserID string `json:"user_id"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Role   string `json:"role"`
	} `json:"user"`
}

// RefreshClaims is minimal claims stored in refresh token.
type RefreshClaims struct {
	Sub   string `json:"sub"`   // subject = user id/email
	Email string `json:"email"` // email
	Role  string `json:"role"`
	jwtStandardClaims
}

// AccessClaims has full access claims
type AccessClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Name  string `json:"name"`
	jwtStandardClaims
}
