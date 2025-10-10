package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/yourname/energy-usage-api/internal/domain"
	"github.com/yourname/energy-usage-api/internal/util"
)

type AuthHandler struct {
	Svc    *domain.Service
	Secret string
}

func NewAuthHandler(s *domain.Service, secret string) *AuthHandler {
	return &AuthHandler{Svc: s, Secret: secret}
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req struct{ Email, Password string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	u, err := h.Svc.Signup(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	_ = json.NewEncoder(w).Encode(u)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct{ Email, Password string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	u, err := h.Svc.VerifyUser(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	jwt, err := util.SignJWT(h.Secret, u.ID, 24*time.Hour)
	if err != nil {
		http.Error(w, "token error", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"token": jwt})
}
