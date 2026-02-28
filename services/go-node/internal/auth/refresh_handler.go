package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// RefreshHandler exposes /api/v1/auth/refresh
type RefreshHandler struct {
	rs     *RefreshService
	domain string // cookie domain (optional)
}

// NewRefreshHandler creates handler
func NewRefreshHandler(rs *RefreshService, cookieDomain string) *RefreshHandler {
	return &RefreshHandler{rs: rs, domain: cookieDomain}
}

// POST /api/v1/auth/refresh
// reads refresh_token cookie, validates, issues new access + refresh and writes cookie
func (h *RefreshHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	c, err := r.Cookie("refresh_token")
	if err != nil || c.Value == "" {
		http.Error(w, `{"error":"missing refresh token"}`, http.StatusUnauthorized)
		return
	}
	accessTok, accessExp, refreshTok, refreshExp, err := h.rs.Rotate(ctx, c.Value)
	if err != nil {
		logrus.Warnf("refresh rotate failed: %v", err)
		http.Error(w, `{"error":"invalid refresh token"}`, http.StatusUnauthorized)
		return
	}
	// write refresh token cookie
	WriteRefreshCookie(w, refreshTok, time.Duration(refreshExp-time.Now().Unix())*time.Second, h.domain)

	// return access token JSON
	resp := map[string]interface{}{
		"access_token": accessTok,
		"token_type":   "bearer",
		"expires_in":   accessExp,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
