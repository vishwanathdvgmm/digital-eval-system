package student

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetResults(w http.ResponseWriter, r *http.Request) {
	usn := r.URL.Query().Get("usn")
	sem := r.URL.Query().Get("semester")
	academicYear := r.URL.Query().Get("academic_year")
	if usn == "" {
		http.Error(w, "missing usn", http.StatusBadRequest)
		return
	}
	res, err := h.svc.FetchResultsWithGPA(r.Context(), usn, sem, academicYear)

	if err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	writeJSON(w, res, http.StatusOK)
}

func (h *Handler) DownloadPDF(w http.ResponseWriter, r *http.Request) {
	usn := r.URL.Query().Get("usn")
	semester := r.URL.Query().Get("semester")
	academicYear := r.URL.Query().Get("academic_year")
	if usn == "" || semester == "" {
		http.Error(w, "missing params", http.StatusBadRequest)
		return
	}

	pdfBytes, err := h.svc.GenerateResultPDF(r.Context(), usn, semester, academicYear)
	if err != nil {
		http.Error(w, "failed to generate pdf: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"result_"+usn+"_"+semester+".pdf\"")
	w.WriteHeader(http.StatusOK)
	w.Write(pdfBytes)
}

func writeJSON(w http.ResponseWriter, v interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func RegisterStudentRoutes(r *mux.Router, svc *Service) {
	h := NewHandler(svc)
	r.HandleFunc("/student/results", h.GetResults).Methods("GET")
	r.HandleFunc("/student/download", h.DownloadPDF).Methods("GET")
}
