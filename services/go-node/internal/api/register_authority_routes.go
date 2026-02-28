package api

import (
	"github.com/gorilla/mux"

	"digital-eval-system/services/go-node/internal/authority"
	"digital-eval-system/services/go-node/internal/core"
)

// RegisterAuthorityRoutes adds authority endpoints if service registered
func RegisterAuthorityRoutes(r *mux.Router, h *Handler, registry *core.ServiceRegistry) {
	if svcIf, ok := registry.Get("authority_service"); ok {
		if svc, ok2 := svcIf.(*authority.Service); ok2 {
			handler := authority.NewHandler(svc)
			r.HandleFunc("/authority/requests/pending", handler.ListPending).Methods("GET")
			r.HandleFunc("/authority/requests/history", handler.ListHistory).Methods("GET")
			r.HandleFunc("/authority/requests/{id}/approve", handler.ApproveRequest).Methods("POST")
			r.HandleFunc("/authority/requests/{id}/reject", handler.RejectRequest).Methods("POST")
		}
	}
	// else: no routes (service not configured)
}
