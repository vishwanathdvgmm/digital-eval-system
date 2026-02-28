package evaluator

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type SubmitHandler struct {
	svc *SubmitService
}

func NewSubmitHandler(svc *SubmitService) *SubmitHandler {
	return &SubmitHandler{svc: svc}
}

func (h *SubmitHandler) Submit(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var in SubmitPayload
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	// ensure numeric course credits default 0
	if in.CourseCredits < 0 {
		in.CourseCredits = 0
	}

	blockHash, err := h.svc.SubmitEvaluation(ctx, in)
	if err != nil {
		logrus.Warnf("submit evaluation failed: %v", err)
		http.Error(w, "submit failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// success response
	resp := map[string]string{
		"block_hash": blockHash,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func RegisterSubmitRoutes(r *mux.Router, svc *SubmitService) {
	r.HandleFunc("/evaluator/submit", NewSubmitHandler(svc).Submit).Methods("POST")
}
