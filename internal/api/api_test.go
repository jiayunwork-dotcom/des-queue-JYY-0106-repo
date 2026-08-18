package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func srv() *Server { return New(Config{Addr: ":0", WebDir: ""}) }

func TestHealth(t *testing.T) {
	w := httptest.NewRecorder()
	srv().Handler().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/health", nil))
	if w.Code != 200 {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestSimulate(t *testing.T) {
	body := `{"lambda":0.5,"mu":1.0,"servers":1,"customers":1000,"seed":42}`
	w := httptest.NewRecorder()
	srv().Handler().ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/simulate", bytes.NewBufferString(body)))
	if w.Code != 200 {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}
	var resp SimResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Served != 1000 {
		t.Fatalf("served = %d", resp.Served)
	}
	if resp.Rho <= 0 {
		t.Fatalf("rho = %f", resp.Rho)
	}
}

func TestTheory(t *testing.T) {
	body := `{"lambda":0.5,"mu":1.0}`
	w := httptest.NewRecorder()
	srv().Handler().ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/theory", bytes.NewBufferString(body)))
	if w.Code != 200 {
		t.Fatalf("status = %d", w.Code)
	}
	var resp TheoryResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Rho != 0.5 {
		t.Fatalf("rho = %f", resp.Rho)
	}
}

func TestSimulateInvalid(t *testing.T) {
	body := `{"lambda":-1,"mu":1.0,"servers":1,"customers":100}`
	w := httptest.NewRecorder()
	srv().Handler().ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/simulate", bytes.NewBufferString(body)))
	if w.Code != 422 {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}
