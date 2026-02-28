package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"digital-eval-system/services/go-node/internal/auth"
)

// RegisterAuthRoutes wires all JWT authentication endpoints.
func RegisterAuthRoutes(r *mux.Router, authSvc *auth.Service, refreshHandler http.Handler, jwtMgr *auth.Manager) {
	base := r.PathPrefix("/api/v1").Subrouter()

	// -------------------------
	// PUBLIC AUTH ENDPOINTS
	// -------------------------
	base.HandleFunc("/auth/login", auth.LoginHandler(authSvc)).Methods("POST")
	base.Handle("/auth/refresh", refreshHandler).Methods("POST")
	base.HandleFunc("/auth/logout", auth.LogoutHandler("")).Methods("POST")

	// -------------------------
	// PROTECTED ROUTES (JWT)
	// -------------------------
	protected := base.NewRoute().Subrouter()
	protected.Use(auth.AuthMiddleware(jwtMgr))

	protected.HandleFunc("/auth/me", func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(auth.ContextKeyUser)
		if u == nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(u)
	}).Methods("GET")
}
