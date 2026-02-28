package auth

import (
	"context"
	"net/http"
	"strings"
)

type ctxKeyType string

const ContextKeyUser ctxKeyType = "auth_user"

func AuthMiddleware(jwtMgr *Manager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, `{"error":"missing auth"}`, http.StatusUnauthorized)
				return
			}
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, `{"error":"invalid auth header"}`, http.StatusUnauthorized)
				return
			}
			token := parts[1]
			claims, err := jwtMgr.ValidateAccessToken(token)
			if err != nil {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}
			// attach user to context
			u := &User{
				ID:     claims.UID,
				UserID: claims.UserID,
				Email:  claims.Email,
				Role:   claims.Role,
				Name:   claims.Name,
			}
			ctx := context.WithValue(r.Context(), ContextKeyUser, u)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole wraps a handler and enforces a role (single role string)
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			v := r.Context().Value(ContextKeyUser)
			if v == nil {
				http.Error(w, `{"error":"unauthenticated"}`, http.StatusUnauthorized)
				return
			}
			u, ok := v.(*User)
			if !ok {
				http.Error(w, `{"error":"unauthenticated"}`, http.StatusUnauthorized)
				return
			}
			if u.Role != role {
				http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (m *Manager) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if h == "" {
			http.Error(w, "missing authorization", http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}
		claims, err := m.ParseAndValidate(parts[1])
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ContextKeyUser, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FromContext returns claims from request context
func FromContext(ctx context.Context) (*User, bool) {
	v := ctx.Value(ContextKeyUser)
	if v == nil {
		return nil, false
	}
	c, ok := v.(*User)
	return c, ok
}
