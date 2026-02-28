package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"digital-eval-system/services/go-node/internal/admin"
	"digital-eval-system/services/go-node/internal/auth"
	"digital-eval-system/services/go-node/internal/authority"
	"digital-eval-system/services/go-node/internal/core"
	"digital-eval-system/services/go-node/internal/evaluator"
	"digital-eval-system/services/go-node/internal/student"
)

func NewRouter(h *Handler) http.Handler {
	r := mux.NewRouter()

	// global middlewares
	r.Use(core.RequestTracingMiddleware)
	// JSON-only except upload
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/v1/examiner/upload" || r.URL.Path == "/api/v1/evaluator/upload" {
				next.ServeHTTP(w, r)
				return
			}
			core.JSONOnlyMiddleware(next).ServeHTTP(w, r)
		})
	})

	r.Use(core.WithTimeoutMiddleware(300 * time.Second))

	apiR := r.PathPrefix("/api/v1").Subrouter()

	// ------------------------------------
	// AUTH ROUTES (JWT v2)
	// ------------------------------------
	av := h.registry.MustGet("auth_service")

	authSvc, ok := av.(*auth.Service)
	if !ok {
		logrus.Fatalf("auth_service wrong type: %T", av)
	}

	authHandler := auth.NewHandler(authSvc)

	// mount /api/v1/auth
	authR := apiR.PathPrefix("/auth").Subrouter()

	// LOGIN (works)
	authR.HandleFunc("/login", authHandler.Login).Methods("POST")

	// REFRESH — using RefreshHandler.Refresh
	refreshSrv := auth.NewRefreshService(authSvc.JWTManager())
	refreshHandler := auth.NewRefreshHandler(refreshSrv, "")

	authR.HandleFunc("/refresh", refreshHandler.Refresh).Methods("POST")

	// LOGOUT
	authR.HandleFunc("/logout", auth.LogoutHandler("")).Methods("POST")

	// base existing routes
	registerValidationRoutes(apiR, h)
	RegisterExaminerRoutes(apiR, h)
	RegisterAuthorityRoutes(apiR, h, h.registry)
	RegisterEvaluatorRoutes(apiR, h, h.registry)

	// -------------------------
	// PHASE 6 NEW ROUTES
	// -------------------------

	// Submit evaluation
	submitSvc := h.registry.MustGet("evaluator_submit_service").(*evaluator.SubmitService)
	RegisterResultsRoutes(apiR, submitSvc)

	// Release results
	releaseSvc := h.registry.MustGet("authority_release_service").(*authority.ReleaseService)
	authority.RegisterReleaseRoutes(apiR, releaseSvc)

	// Student result access (correct mounting under /api/v1)
	studentSvc := h.registry.MustGet("student_service").(*student.Service)
	studentHandler := student.NewHandler(studentSvc)

	// MOUNT inside /api/v1
	apiR.HandleFunc("/student/results", studentHandler.GetResults).Methods("GET")
	apiR.HandleFunc("/student/download", studentHandler.DownloadPDF).Methods("GET")

	// ADMIN ROUTES
	if val, ok := h.registry.Get("admin_service"); ok {
		if adminSvc, ok := val.(*admin.Service); ok {
			RegisterAdminRoutes(apiR, adminSvc)
		}
	}

	// chain/block internal routes
	apiR.HandleFunc("/blocks", h.HandlePostBlock).Methods("POST")
	apiR.HandleFunc("/blocks/{hash}", h.HandleGetBlock).Methods("GET")
	apiR.HandleFunc("/chain/height", h.HandleGetHead).Methods("GET")
	apiR.HandleFunc("/chain/verify", h.HandleVerifyChain).Methods("GET")

	// health
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	return CORSMiddleware(r)
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // In production, replace * with specific origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Pass to next handler
		next.ServeHTTP(w, r)
	})
}
