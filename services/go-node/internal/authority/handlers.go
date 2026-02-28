package authority

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Handler wraps service and exposes HTTP endpoints.
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// GET /api/v1/authority/requests/pending
func (h *Handler) ListPending(w http.ResponseWriter, r *http.Request) {
	rows, err := h.svc.ListPending(r.Context())
	if err != nil {
		http.Error(w, "failed to load", http.StatusInternalServerError)
		return
	}
	writeJSON(w, rows, http.StatusOK)
}

// GET /api/v1/authority/requests/history
func (h *Handler) ListHistory(w http.ResponseWriter, r *http.Request) {
	rows, err := h.svc.ListHistory(r.Context())
	if err != nil {
		http.Error(w, "failed to load history", http.StatusInternalServerError)
		return
	}
	writeJSON(w, rows, http.StatusOK)
}

// POST /api/v1/authority/requests/{id}/approve
func (h *Handler) ApproveRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	// payload to get evaluator and assign count
	var p ApprovePayload
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		// allow missing payload, default assign 5
		p.AssignNum = 5
	}
	// For simplicity fetch the request row to get evaluatorID, course, semester
	// (We will call ListPending and filter; production should have direct query)
	rows, err := h.svc.ListPending(r.Context())
	if err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	var target *RequestRow
	for _, rrow := range rows {
		if rrow.ID == id {
			target = &RequestRow{
				ID:           rrow.ID,
				EvaluatorID:  rrow.EvaluatorID,
				CourseID:     rrow.CourseID,
				Semester:     rrow.Semester,
				AcademicYear: rrow.AcademicYear,
			}
			break
		}
	}
	if target == nil {
		http.Error(w, "request not found", http.StatusNotFound)
		return
	}

	assigned, err := h.svc.ApproveRequest(r.Context(), id, target.EvaluatorID, target.CourseID, target.Semester, target.AcademicYear, p.AssignNum)
	if err != nil {
		http.Error(w, "approve failed", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{"assigned": assigned}, http.StatusOK)
}

// POST /api/v1/authority/requests/{id}/reject
func (h *Handler) RejectRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.svc.RejectRequest(r.Context(), id); err != nil {
		http.Error(w, "reject failed", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{"status": "rejected"}, http.StatusOK)
}
