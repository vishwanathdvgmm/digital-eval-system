package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns bcrypt hash for a plaintext password.
func HashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}

// ComparePassword returns nil when password matches hashed.
func ComparePassword(hash, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}
