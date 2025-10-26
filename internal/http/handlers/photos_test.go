package handlers

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/yourname/energy-usage-api/internal/domain"
	"github.com/yourname/energy-usage-api/internal/pkg/userctx"
)

// Satisfy the *domain.Service interface for PhotosHandler in test.
type fakeService struct{}

func (f *fakeService) SaveMeterPhoto(ctx context.Context, photo *domain.MeterPhoto) error {
	return nil
}

func TestPhotosHandler_Upload(t *testing.T) {
	tmpDir := t.TempDir()
	handler := &PhotosHandler{
		Svc: &fakeService{},
		Dir: tmpDir,
	}

	// Prepare a fake file upload
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("photo", "test.jpg")
	part.Write([]byte("fake image data"))
	writer.Close()

	req := httptest.NewRequest("POST", "/v1/meters/meter123/photos/faulty", body)
	ctx := context.WithValue(req.Context(), userctx.UserIDKey, "user123")
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler.Upload(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rr.Code)
	}

	files, _ := os.ReadDir(tmpDir)
	if len(files) == 0 {
		t.Fatal("expected a file to be saved")
	}
}
