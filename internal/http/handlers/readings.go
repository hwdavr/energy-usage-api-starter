package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/yourname/energy-usage-api/internal/domain"
)

type ReadingsHandler struct{ Svc *domain.Service }

func (h *ReadingsHandler) Add(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ReadingKwh float64   `json:"reading_kwh"`
		ReadingAt  time.Time `json:"reading_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest); return
	}
	meterID := chi.URLParam(r, "meterID")
	rd, err := h.Svc.AddReading(r.Context(), meterID, req.ReadingKwh, req.ReadingAt)
	if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
	_ = json.NewEncoder(w).Encode(rd)
}

func (h *ReadingsHandler) Usage(w http.ResponseWriter, r *http.Request) {
	meterID := chi.URLParam(r, "meterID")
	fromStr := r.URL.Query().Get("from")
	toStr   := r.URL.Query().Get("to")
	from, err := time.Parse(time.RFC3339, fromStr); if err != nil { http.Error(w, "bad from", http.StatusBadRequest); return }
	to,   err := time.Parse(time.RFC3339, toStr);   if err != nil { http.Error(w, "bad to", http.StatusBadRequest); return }
	total, err := h.Svc.Usage(r.Context(), meterID, from, to)
	if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
	_ = json.NewEncoder(w).Encode(map[string]any{"meter_id": meterID, "kwh": total, "from": from, "to": to})
}
