package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/yourname/energy-usage-api/internal/domain"
	h "github.com/yourname/energy-usage-api/internal/http/handlers"
)

func TestSignupLogin(t *testing.T) {
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL not set; skipping integration-ish test")
	}
	db, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	repo := domain.NewRepository(db)
	svc := domain.NewService(repo)
	ah := h.NewAuthHandler(svc, "testsecret")

	// signup
	req := httptest.NewRequest("POST", "/v1/auth/signup", bytes.NewBufferString(`{"Email":"t@e.com","Password":"x"}`))
	w := httptest.NewRecorder()
	ah.Signup(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("signup code=%d", w.Code)
	}
	// login
	req = httptest.NewRequest("POST", "/v1/auth/login", bytes.NewBufferString(`{"Email":"t@e.com","Password":"x"}`))
	w = httptest.NewRecorder()
	ah.Login(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login code=%d", w.Code)
	}
}
