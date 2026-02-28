package evaluator

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Handler exposes evaluator endpoints
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateRequest(w http.ResponseWriter, r *http.Request) {
	var payload RequestCreate
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	id, err := h.svc.CreateRequest(r.Context(), payload.EvaluatorID, payload.CourseID, payload.Semester, payload.AcademicYear, payload.Description)
	if err != nil {
		logrus.Warnf("create request failed: %v", err)
		if strings.Contains(err.Error(), "uq_eval_request") {
			http.Error(w, "request already exists", http.StatusConflict)
			return
		}
		return
	}
	writeJSON(w, map[string]interface{}{"request_id": id}, http.StatusCreated)
}

// GET /api/v1/evaluator/requests/history?evaluator_id=...
func (h *Handler) ListRequests(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	eid := q.Get("evaluator_id")
	if eid == "" {
		http.Error(w, "missing evaluator_id", http.StatusBadRequest)
		return
	}
	rows, err := h.svc.ListRequests(r.Context(), eid)
	if err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	writeJSON(w, rows, http.StatusOK)
}

// GET /api/v1/evaluator/assigned?evaluator_id=...
func (h *Handler) ListAssigned(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	eid := q.Get("evaluator_id")
	if eid == "" {
		http.Error(w, "missing evaluator_id", http.StatusBadRequest)
		return
	}
	rows, err := h.svc.ListAssigned(r.Context(), eid)
	if err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	writeJSON(w, rows, http.StatusOK)
}

// GET /api/v1/evaluator/script/{script_id}
func (h *Handler) GetScript(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sid := vars["script_id"]
	if sid == "" {
		http.Error(w, "missing script_id", http.StatusBadRequest)
		return
	}
	meta, err := h.svc.GetScriptMetadata(r.Context(), sid)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, meta, http.StatusOK)
}

// POST /api/v1/evaluator/submit
func (h *Handler) Submit(w http.ResponseWriter, r *http.Request) {
	var payload SubmitPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	blockHash, err := h.svc.SubmitEvaluation(r.Context(), payload)
	if err != nil {
		http.Error(w, "submit failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"block_hash": blockHash}, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, v interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
