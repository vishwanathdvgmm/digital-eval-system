package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"digital-eval-system/services/go-node/internal/examiner"
	"digital-eval-system/services/go-node/internal/rootdir"
)

// RegisterExaminerRoutes wires the upload handler at /api/v1/examiner/upload
func RegisterExaminerRoutes(r *mux.Router, h *Handler) {
	// build upload handler using service from registry if present
	var uploadHandler http.Handler

	// attempt to get service from registry (we created service in main.go below)
	if svcIf, ok := h.registry.Get("examiner_upload_service"); ok {
		if svc, ok2 := svcIf.(*examiner.Service); ok2 {
			uploadHandler = examiner.NewUploadHandler(svc, rootdir.Resolve("data/uploads"), 24*3600)
		}
	}

	// fallback: if not registered, do nothing (route not created)
	if uploadHandler != nil {
		r.Handle("/examiner/upload", uploadHandler).Methods("POST")
	}
}
