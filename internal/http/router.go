package httpserver

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/AyomiCoder/dstr-sys.git/internal/config"
	"github.com/AyomiCoder/dstr-sys.git/internal/notification"
)

const maxNotificationBodyBytes int64 = 1 << 20 // 1MB

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

func NewRouter(cfg config.Config, repo notification.Repository) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler(cfg))
	mux.HandleFunc("POST /notifications", createNotificationHandler(repo))
	mux.HandleFunc("GET /notifications", listNotificationsHandler(repo))

	return mux
}

func createNotificationHandler(repo notification.Repository) http.HandlerFunc {
	type createNotificationRequest struct {
		Type      string `json:"type"`
		Recipient string `json:"recipient"`
		Message   string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if contentType != "" && !strings.HasPrefix(contentType, "application/json") {
			writeJSONError(w, http.StatusUnsupportedMediaType, "content type must be application/json")
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxNotificationBodyBytes)
		defer r.Body.Close()

		var req createNotificationRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&req); err != nil {
			var syntaxErr *json.SyntaxError
			switch {
			case errors.As(err, &syntaxErr):
				writeJSONError(w, http.StatusBadRequest, "malformed json body")
			case errors.Is(err, io.EOF):
				writeJSONError(w, http.StatusBadRequest, "request body is required")
			default:
				writeJSONError(w, http.StatusBadRequest, "invalid request body")
			}
			return
		}

		var extra json.RawMessage
		if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
			writeJSONError(w, http.StatusBadRequest, "request body must contain a single json object")
			return
		}

		req.Type = strings.TrimSpace(req.Type)
		req.Recipient = strings.TrimSpace(req.Recipient)
		req.Message = strings.TrimSpace(req.Message)

		if req.Type == "" || req.Recipient == "" || req.Message == "" {
			writeJSONError(w, http.StatusBadRequest, "type, recipient, and message are required")
			return
		}

		created, err := repo.Create(r.Context(), notification.CreateInput{
			Type:      req.Type,
			Recipient: req.Recipient,
			Message:   req.Message,
		})
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to create notification")
			return
		}

		writeJSON(w, http.StatusCreated, created)
	}
}

func listNotificationsHandler(repo notification.Repository) http.HandlerFunc {
	type listNotificationsResponse struct {
		Notifications []notification.Notification `json:"notifications"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		items, err := repo.List(r.Context())
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch notifications")
			return
		}

		writeJSON(w, http.StatusOK, listNotificationsResponse{Notifications: items})
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	writeJSON(w, statusCode, errorResponse{Error: message})
}
