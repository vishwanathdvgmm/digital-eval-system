package api

import (
	"github.com/gorilla/mux"

	"digital-eval-system/services/go-node/internal/evaluator"
)

// call this from routes.go with proper services
func RegisterResultsRoutes(r *mux.Router, submitSvc *evaluator.SubmitService) {
	if submitSvc != nil {
		evaluator.RegisterSubmitRoutes(r, submitSvc)
	}
}
