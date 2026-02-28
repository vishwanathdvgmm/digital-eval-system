package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"digital-eval-system/services/go-node/internal/admin"
)

func RegisterAdminRoutes(r *mux.Router, svc *admin.Service) {
	s := r.PathPrefix("/admin").Subrouter()

	// POST /services/{name}/{action}
	s.HandleFunc("/services/{name}/{action}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		action := vars["action"]

		var err error
		switch action {
		case "start":
			err = svc.Start(name)
		case "stop":
			err = svc.Stop(name)
		case "restart":
			err = svc.Restart(name)
		default:
			http.Error(w, "invalid action", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": action + " executed"})
	}).Methods("POST")

	// GET /services/{name}/logs
	s.HandleFunc("/services/{name}/logs", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		logs, err := svc.GetLogs(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"logs": logs})
	}).Methods("GET")
}
