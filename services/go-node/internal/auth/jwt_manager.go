package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTManager handles RSA256 tokens
type Manager struct {
	privKey *rsa.PrivateKey
	pubKey  *rsa.PublicKey

	AccessTTL  time.Duration
	RefreshTTL time.Duration
	Issuer     string
	Audience   string
}

// NewJWTManager loads keys from paths and returns manager
func NewManagerFromFiles(privPath, pubPath string, accessTTLMinutes int, refreshTTLDays int, issuer, audience string) (*Manager, error) {
	privBytes, err := os.ReadFile(privPath)
	if err != nil {
		return nil, fmt.Errorf("read priv: %w", err)
	}
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		return nil, fmt.Errorf("parse priv: %w", err)
	}

	pubBytes, err := os.ReadFile(pubPath)
	if err != nil {
		return nil, fmt.Errorf("read pub: %w", err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		return nil, fmt.Errorf("parse pub: %w", err)
	}

	return &Manager{
		privKey:    privKey,
		pubKey:     pubKey,
		AccessTTL:  time.Duration(accessTTLMinutes) * time.Minute,
		RefreshTTL: time.Duration(refreshTTLDays) * 24 * time.Hour,
		Issuer:     issuer,
		Audience:   audience,
	}, nil
}

var ErrInvalidToken = errors.New("invalid token")

type Claims struct {
	UID    int64  `json:"uid"`
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

func (m *Manager) ParseAndValidate(tokenStr string) (*Claims, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}))
	tok, err := parser.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return m.pubKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := tok.Claims.(*Claims); ok && tok.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

// GenerateAccessToken creates RSA-signed JWT
func (m *Manager) GenerateAccessToken(u *User) (string, error) {
	// set TTL from jm.config or default
	now := time.Now().UTC()
	claims := &Claims{
		UID:    u.ID,
		UserID: u.UserID,
		Email:  u.Email,
		Role:   u.Role,
		Name:   u.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.Issuer,
			Audience:  jwt.ClaimStrings{m.Audience},
			Subject:   u.UserID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.AccessTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(m.privKey)
}

func (m *Manager) GenerateRefreshToken(u *User) (string, error) {
	now := time.Now().UTC()
	claims := &Claims{
		UID:    u.ID,
		UserID: u.UserID,
		Email:  u.Email,
		Role:   u.Role,
		Name:   u.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.Issuer,
			Audience:  jwt.ClaimStrings{m.Audience},
			Subject:   u.UserID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.RefreshTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(m.privKey)
}

// ValidateAccessToken verifies signature + expiry and returns claims
func (m *Manager) ValidateAccessToken(tokenStr string) (*Claims, error) {
	parser := jwt.NewParser()
	claims := &Claims{}
	_, err := parser.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodRS256 {
			return nil, errors.New("unexpected signing method")
		}
		return m.pubKey, nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (m *Manager) ValidateRefreshToken(tokenStr string) (*Claims, error) {
	return m.ValidateAccessToken(tokenStr)
}
