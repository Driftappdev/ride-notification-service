package httpadapter

import (
	"encoding/json"
	"net/http"
	"strings"

	"dift_backend_go/notification-service/internal/service"
)

type AdminControlHandler struct {
	svc          *service.AdminControlService
	sharedSecret string
}

func NewAdminControlHandler(svc *service.AdminControlService, sharedSecret string) *AdminControlHandler {
	return &AdminControlHandler{svc: svc, sharedSecret: sharedSecret}
}

func (h *AdminControlHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/internal/admin/control", h.handleControl)
}

func (h *AdminControlHandler) handleControl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !h.authorized(r) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "unauthorized"})
		return
	}

	var req service.AdminControlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "invalid_json"})
		return
	}

	resp, err := h.svc.Execute(r.Context(), r.Header.Get("Idempotency-Key"), req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *AdminControlHandler) authorized(r *http.Request) bool {
	if strings.TrimSpace(h.sharedSecret) == "" {
		return true
	}
	return r.Header.Get("X-Admin-Secret") == h.sharedSecret
}
