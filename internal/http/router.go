package httpserver

import (
	"encoding/json"
	"net/http"

	"dstr_sys/internal/config"
)

func NewRouter(cfg config.Config) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler(cfg))

	return mux
}

func healthHandler(cfg config.Config) http.HandlerFunc {
	type healthResponse struct {
		Status  string `json:"status"`
		Service string `json:"service"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(healthResponse{
			Status:  "ok",
			Service: cfg.ServiceName,
		}); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
