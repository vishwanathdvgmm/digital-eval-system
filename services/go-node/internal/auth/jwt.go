package auth

// Claims returned by JWT
type TokenClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
	jwtStandardClaims
}

type jwtStandardClaims struct {
	ExpiresAt int64 `json:"exp"`
	IssuedAt  int64 `json:"iat"`
	NotBefore int64 `json:"nbf"`
	// add jti if needed
}
