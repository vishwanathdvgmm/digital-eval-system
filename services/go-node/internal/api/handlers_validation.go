package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"digital-eval-system/services/go-node/internal/pybridge"
)

// ValidationRequest represents the incoming request to trigger extraction + validation
type ValidationRequest struct {
	FilePath string `json:"file_path"`
}

// ValidationResponse returns combined extraction + validation result
type ValidationResponse struct {
	Extract *pybridge.ExtractResponse  `json:"extract,omitempty"`
	Valid   *pybridge.ValidateResponse `json:"validate,omitempty"`
	Error   string                     `json:"error,omitempty"`
}

// HandleValidateFile triggers Python extraction then validation and returns results.
// POST /api/v1/validate
func (h *Handler) HandleValidateFile(w http.ResponseWriter, r *http.Request) {
	var req ValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.FilePath == "" {
		httpError(w, "missing file_path", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 300*time.Second)
	defer cancel()

	// 1) call python extract
	exOut, err := h.pyClient().Extract(ctx, req.FilePath)
	if err != nil {
		logrus.Warnf("py extract failed: %v", err)
		writeJSON(w, &ValidationResponse{Extract: exOut, Error: err.Error()}, http.StatusBadGateway)
		return
	}

	// If extraction returned status != success, return early
	if exOut == nil || exOut.Status != "success" {
		writeJSON(w, &ValidationResponse{Extract: exOut, Error: "extraction failed"}, http.StatusBadGateway)
		return
	}

	// 2) call python validate with metadata
	valOut, err := h.pyClient().Validate(ctx, exOut.Metadata)
	if err != nil {
		logrus.Warnf("py validate failed: %v", err)
		writeJSON(w, &ValidationResponse{Extract: exOut, Valid: valOut, Error: err.Error()}, http.StatusBadGateway)
		return
	}

	writeJSON(w, &ValidationResponse{Extract: exOut, Valid: valOut}, http.StatusOK)
}

// Add route registration helper for validation
func registerValidationRoutes(r *mux.Router, h *Handler) {
	r.HandleFunc("/validate", h.HandleValidateFile).Methods("POST")
}
