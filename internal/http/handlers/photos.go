// Package handlers provides HTTP handlers for meter photo uploads.
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/yourname/energy-usage-api/internal/domain"
	"github.com/yourname/energy-usage-api/internal/pkg/userctx"
)

type PhotoType string

const (
	PhotoFaulty PhotoType = "faulty"
	PhotoFixed  PhotoType = "fixed"
)

type PhotosHandler struct {
	// Svc is any type that can save meter photo metadata. Use an interface
	// so tests can inject a fake implementation without depending on the
	// concrete domain.Service type.
	Svc interface {
		SaveMeterPhoto(ctx context.Context, photo *domain.MeterPhoto) error
	}
	// Dir is the directory to store uploaded photos (local disk for now)
	Dir string
}

func (h *PhotosHandler) Upload(w http.ResponseWriter, r *http.Request) {
	meterID := chi.URLParam(r, "meterID")
	photoType := chi.URLParam(r, "type") // "faulty" or "fixed"
	val := r.Context().Value(userctx.UserIDKey)
	userID, ok := val.(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	file, handler, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "missing photo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Determine directory: handler.Dir -> env -> default
	dir := h.Dir
	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		http.Error(w, "failed to save photo", http.StatusInternalServerError)
		return
	}

	// Save file to disk (local)
	ext := filepath.Ext(handler.Filename)
	filename := fmt.Sprintf("%s_%s_%d%s", meterID, photoType, time.Now().UnixNano(), ext)
	savePath := filepath.Join(dir, filename)
	fmt.Printf("Save file to path: %s\n", savePath)
	out, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "failed to save photo", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, "failed to write photo", http.StatusInternalServerError)
		return
	}

	// Store metadata
	meta := domain.MeterPhoto{
		MeterID:   meterID,
		UserID:    userID,
		Type:      string(photoType),
		Path:      filename,
		CreatedAt: time.Now(),
	}
	if err := h.Svc.SaveMeterPhoto(r.Context(), &meta); err != nil {
		http.Error(w, "failed to save metadata", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(meta)
}
