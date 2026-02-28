package auth

import (
	"encoding/json"
	"net/http"
	"time"
)

// LoginRequest used for /auth/login
type LoginRequest struct {
	Login    string `json:"login"`    // email
	Password string `json:"password"` // plaintext
}

// LoginResponse returned on successful login
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	User        *User  `json:"user"`
}

// RefreshService provides GenerateAccess from refresh token
type RefreshService struct {
	jwt *Manager
}

func NewRefreshService(jwt *Manager) *RefreshService {
	return &RefreshService{jwt: jwt}
}

// LoginHandler : POST /api/v1/auth/login
func LoginHandler(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		user, err := svc.Authenticate(ctx, in.Login, in.Password)
		if err != nil {
			http.Error(w, `{"error":"invalid login"}`, http.StatusUnauthorized)
			return
		}

		access, err := svc.jwt.GenerateAccessToken(user)
		if err != nil {
			http.Error(w, "token generation failed", http.StatusInternalServerError)
			return
		}
		refresh, err := svc.jwt.GenerateRefreshToken(user)
		if err != nil {
			http.Error(w, "refresh token generation failed", http.StatusInternalServerError)
			return
		}

		// set refresh cookie (HttpOnly)
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refresh,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			// Secure: true,  // enable if using HTTPS
			Expires: time.Now().Add(svc.jwt.RefreshTTL),
		})

		resp := LoginResponse{
			AccessToken: access,
			TokenType:   "bearer",
			User: &User{
				ID:     user.ID,
				UserID: user.UserID,
				Email:  user.Email,
				Role:   user.Role,
				Name:   user.Name,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

// Logout handler clears refresh cookie and (optionally) revoke refresh token server-side
func LogoutHandler(cookieName string) http.HandlerFunc {
	if cookieName == "" {
		cookieName = "refresh_token"
	}
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
			Expires:  time.Unix(0, 0),
		})
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}
}

// Handler aggregates auth endpoints
type Handler struct {
	svc            *Service
	refreshHandler *RefreshHandler
}

func NewHandler(svc *Service) *Handler {
	rs := NewRefreshService(svc.jwt)
	// assuming empty domain for now, or could be passed in if needed
	rh := NewRefreshHandler(rs, "")
	return &Handler{
		svc:            svc,
		refreshHandler: rh,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	LoginHandler(h.svc)(w, r)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	h.refreshHandler.Refresh(w, r)
}
