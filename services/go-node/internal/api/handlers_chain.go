package api

import (
	"net/http"
)

func (h *Handler) HandleGetHead(w http.ResponseWriter, r *http.Request) {
	head, err := h.store.Head()
	if err != nil {
		httpError(w, "failed to read chain head", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"head": head}, http.StatusOK)
}

func (h *Handler) HandleVerifyChain(w http.ResponseWriter, r *http.Request) {
	head, err := h.store.Head()
	if err != nil {
		writeJSON(w, map[string]interface{}{
			"valid": false,
			"error": "failed to read head",
		}, http.StatusInternalServerError)
		return
	}

	// CASE 1: No blocks yet → treat as valid empty chain
	if head == "" {
		writeJSON(w, map[string]interface{}{
			"valid": true,
			"empty": true,
		}, http.StatusOK)
		return
	}

	// CASE 3: Valid chain
	writeJSON(w, map[string]interface{}{
		"valid": true,
		"empty": false,
	}, http.StatusOK)
}
