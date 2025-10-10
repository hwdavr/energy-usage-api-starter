package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yourname/energy-usage-api/internal/domain"
	"github.com/yourname/energy-usage-api/internal/pkg/userctx"
)

type MetersHandler struct{ Svc *domain.Service }

func (h *MetersHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct{ Label string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	val := r.Context().Value(userctx.UserIDKey)
	userID, ok := val.(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	m, err := h.Svc.CreateMeter(r.Context(), userID, req.Label)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(m)
}

func (h *MetersHandler) List(w http.ResponseWriter, r *http.Request) {
	val := r.Context().Value(userctx.UserIDKey)
	userID, ok := val.(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	ms, err := h.Svc.ListMeters(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(ms)
}
