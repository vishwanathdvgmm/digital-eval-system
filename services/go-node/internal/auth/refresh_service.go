package auth

import (
	"context"
	"time"
)

// Rotate validates incoming refresh token and issues new pair.
func (rs *RefreshService) Rotate(ctx context.Context, tokenStr string) (accessToken string, accessExp int64, refreshToken string, refreshExp int64, err error) {
	claims, err := rs.jwt.ValidateRefreshToken(tokenStr)
	if err != nil {
		return "", 0, "", 0, err
	}
	// extract required fields
	email := claims.Email
	role := claims.Role
	name := claims.Name
	userID := claims.UserID

	// reconstruct user for generation
	u := &User{
		UserID: userID,
		Email:  email,
		Role:   role,
		Name:   name,
		// ID is likely needed if it's used in claims, but Claims struct has UID
		ID: claims.UID,
	}

	// issue new tokens
	accessToken, err = rs.jwt.GenerateAccessToken(u)
	if err != nil {
		return "", 0, "", 0, err
	}
	// Calculate expiration based on TTL
	accessExp = time.Now().Add(rs.jwt.AccessTTL).Unix()

	refreshToken, err = rs.jwt.GenerateRefreshToken(u)
	if err != nil {
		return "", 0, "", 0, err
	}
	// Calculate expiration based on TTL
	refreshExp = time.Now().Add(rs.jwt.RefreshTTL).Unix()

	return accessToken, accessExp, refreshToken, refreshExp, nil
}

// ValidateRefresh is a helper to just validate and return claims
func (rs *RefreshService) ValidateRefresh(tokenStr string) (map[string]interface{}, error) {
	claims, err := rs.jwt.ValidateRefreshToken(tokenStr)
	if err != nil {
		return nil, err
	}
	out := map[string]interface{}{
		"email":   claims.Email,
		"role":    claims.Role,
		"user_id": claims.UserID,
		"uid":     claims.UID,
		"name":    claims.Name,
	}
	return out, nil
}
