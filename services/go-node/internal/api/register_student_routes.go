package api

import (
	"github.com/gorilla/mux"

	"digital-eval-system/services/go-node/internal/student"
)

// RegisterStudentRoutes wires student endpoints
func RegisterStudentRoutes(r *mux.Router, studentSvc *student.Service) {
	if studentSvc != nil {
		student.RegisterStudentRoutes(r, studentSvc)
	}
}
