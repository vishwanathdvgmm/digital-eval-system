package auth

import (
	"net/http"
	"time"
)

// WriteRefreshCookie writes HttpOnly secure cookie for refresh token.
// Path and domain can be adjusted by caller; defaults set to "/".
func WriteRefreshCookie(w http.ResponseWriter, token string, ttl time.Duration, domain string) {
	c := &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // you run http now; set true for HTTPS production
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(ttl.Seconds()),
	}
	if domain != "" {
		c.Domain = domain
	}
	http.SetCookie(w, c)
}

// ClearRefreshCookie deletes cookie.
func ClearRefreshCookie(w http.ResponseWriter, domain string) {
	c := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	}
	if domain != "" {
		c.Domain = domain
	}
	http.SetCookie(w, c)
}
