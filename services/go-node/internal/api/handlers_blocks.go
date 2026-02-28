package api

import (
	"crypto/rsa"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"digital-eval-system/services/go-node/internal/block"
)

func (h *Handler) HandlePostBlock(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		httpError(w, "failed to read body", http.StatusBadRequest)
		return
	}

	var b block.Block
	if err := json.Unmarshal(body, &b); err != nil {
		httpError(w, "invalid block json", http.StatusBadRequest)
		return
	}

	if len(b.Header.Signature) == 0 {
		httpError(w, "header signature missing", http.StatusBadRequest)
		return
	}

	// Required field validation
	for _, tx := range b.Transactions {
		if !block.ValidateTransaction(&tx) {
			httpError(w, "invalid transaction payload", http.StatusBadRequest)
			return
		}
	}

	// Append block
	hash, err := h.chain.AppendBlock(&b)
	if err != nil {
		httpError(w, "failed to append block", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"block_hash": hash}, http.StatusCreated)
}

func (h *Handler) HandleGetBlock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	if hash == "" {
		httpError(w, "missing block hash", http.StatusBadRequest)
		return
	}

	b, err := h.store.GetBlock(hash)
	if err != nil {
		httpError(w, "block not found", http.StatusNotFound)
		return
	}

	writeJSON(w, b, http.StatusOK)
}

var _ = (*rsa.PublicKey)(nil)

// helpers
func writeJSON(w http.ResponseWriter, v interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func httpError(w http.ResponseWriter, msg string, code int) {
	logrus.Warnf("http error %d: %s", code, msg)
	writeJSON(w, map[string]string{"error": msg}, code)
}
