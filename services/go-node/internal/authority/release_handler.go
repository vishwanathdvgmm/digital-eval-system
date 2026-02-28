package authority

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type ReleaseHandler struct {
	svc *ReleaseService
}

func NewReleaseHandler(svc *ReleaseService) *ReleaseHandler {
	return &ReleaseHandler{svc: svc}
}

// POST /api/v1/authority/results/release
// payload: { "semester": "5", "academic_year": "2024-2025", "released_by": "authority_1" }
func (h *ReleaseHandler) Release(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Semester     string `json:"semester"`
		AcademicYear string `json:"academic_year"`
		ReleasedBy   string `json:"released_by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if payload.Semester == "" || payload.AcademicYear == "" || payload.ReleasedBy == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}

	blockHash, err := h.svc.ReleaseResults(r.Context(), payload.Semester, payload.AcademicYear, payload.ReleasedBy)
	if err != nil {
		http.Error(w, "release failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"block_hash": blockHash}, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, v interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func RegisterReleaseRoutes(r *mux.Router, svc *ReleaseService) {
	r.HandleFunc("/authority/results/release", NewReleaseHandler(svc).Release).Methods("POST")
}
