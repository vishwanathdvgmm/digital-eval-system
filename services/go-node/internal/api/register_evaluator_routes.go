package api

import (
	"github.com/gorilla/mux"

	"digital-eval-system/services/go-node/internal/core"
	"digital-eval-system/services/go-node/internal/evaluator"
)

// RegisterEvaluatorRoutes wires evaluator endpoints if registered
func RegisterEvaluatorRoutes(r *mux.Router, h *Handler, registry *core.ServiceRegistry) {

	svcIf, ok := registry.Get("evaluator_service")
	if !ok {
		return
	}

	svc, ok := svcIf.(*evaluator.Service)
	if !ok {
		return
	}

	hd := evaluator.NewHandler(svc)

	// ============================
	// Evaluator endpoints
	// ============================

	// Evaluator creates request to authority
	r.HandleFunc("/evaluator/requests", hd.CreateRequest).Methods("POST")

	// Evaluator fetches request history
	r.HandleFunc("/evaluator/requests/history", hd.ListRequests).Methods("GET")

	// Evaluator fetches assigned scripts
	r.HandleFunc("/evaluator/assigned", hd.ListAssigned).Methods("GET")

	// Evaluator fetches an assigned script
	r.HandleFunc("/evaluator/script/{script_id}", hd.GetScript).Methods("GET")

	// Evaluator submits evaluation (Marks)
	if submitSvcIf, ok := registry.Get("evaluator_submit_service"); ok {
		submitSvc := submitSvcIf.(*evaluator.SubmitService)
		submitHandler := evaluator.NewSubmitHandler(submitSvc)
		r.HandleFunc("/evaluator/submit", submitHandler.Submit).Methods("POST")
	}

	// Evaluator uploads evaluated script (File) - NEW
	if uploadSvcIf, ok := registry.Get("evaluator_upload_service"); ok {
		uploadSvc := uploadSvcIf.(*evaluator.UploadService)
		uploadHandler := evaluator.NewUploadHandler(uploadSvc, "")
		r.Handle("/evaluator/upload", uploadHandler).Methods("POST")
	}
}
