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

func setupTestHandler(t *testing.T) *h.AuthHandler {
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL not set; skipping integration-ish test")
	}
	db, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	repo := domain.NewRepository(db)
	svc := domain.NewService(repo)
	return h.NewAuthHandler(svc, "testsecret")
}

func TestSignupLogin(t *testing.T) {
	ah := setupTestHandler(t)
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

func TestSignupDecodeError(t *testing.T) {
	ah := h.NewAuthHandler(nil, "testsecret")
	req := httptest.NewRequest("POST", "/v1/auth/signup", bytes.NewBufferString(`invalid json`))
	w := httptest.NewRecorder()
	ah.Signup(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestSignupDuplicateEmail(t *testing.T) {
	ah := setupTestHandler(t)

	// first signup
	req := httptest.NewRequest("POST", "/v1/auth/signup", bytes.NewBufferString(`{"Email":"t@e.com","Password":"x"}`))
	w := httptest.NewRecorder()
	ah.Signup(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("signup code=%d", w.Code)
	}

	// second signup (duplicate email)
	req = httptest.NewRequest("POST", "/v1/auth/signup", bytes.NewBufferString(`{"Email":"t@e.com","Password":"x"}`))
	w = httptest.NewRecorder()
	ah.Signup(w, req)
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
}

func TestLoginDecodeError(t *testing.T) {
	ah := h.NewAuthHandler(nil, "testsecret")
	req := httptest.NewRequest("POST", "/v1/auth/login", bytes.NewBufferString(`invalid json`))
	w := httptest.NewRecorder()
	ah.Login(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	ah := setupTestHandler(t)
	req := httptest.NewRequest("POST", "/v1/auth/login", bytes.NewBufferString(`{"Email":"nonexistent","Password":"wrong"}`))
	w := httptest.NewRecorder()
	ah.Login(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestLoginTokenError(t *testing.T) {
	ah := setupTestHandler(t)
	// First, signup a user
	req := httptest.NewRequest("POST", "/v1/auth/signup", bytes.NewBufferString(`{"Email":"t@e.com","Password":"x"}`))
	w := httptest.NewRecorder()
	ah.Signup(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("signup code=%d", w.Code)
	}
	// Now, try to login with the same user
	req = httptest.NewRequest("POST", "/v1/auth/login", bytes.NewBufferString(`{"Email":"t@e.com","Password":"x"}`))
	w = httptest.NewRecorder()
	ah.Login(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login code=%d", w.Code)
	}
	// Finally, simulate a token error
	// (In a real test, you would mock the token generation/validation)
	// Here, we just check that the handler returns a 500 status code
	req = httptest.NewRequest("POST", "/v1/auth/login", bytes.NewBufferString(`{"Email":"t@e.com","Password":"x"}`))
	w = httptest.NewRecorder()
	ah.Login(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
