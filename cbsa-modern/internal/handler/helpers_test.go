package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	JSON(w, http.StatusOK, map[string]string{"key": "value"})

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
	body := strings.TrimSpace(w.Body.String())
	if body != `{"key":"value"}` {
		t.Errorf("unexpected body: %s", body)
	}
}

func TestErrorJSON(t *testing.T) {
	w := httptest.NewRecorder()
	ErrorJSON(w, http.StatusBadRequest, "bad input")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
	body := strings.TrimSpace(w.Body.String())
	if !strings.Contains(body, "bad input") {
		t.Errorf("expected body to contain 'bad input', got %s", body)
	}
}

func TestQueryInt(t *testing.T) {
	tests := []struct {
		name       string
		queryParam string
		key        string
		defaultVal int
		expected   int
	}{
		{"valid", "?limit=10", "limit", 20, 10},
		{"missing", "", "limit", 20, 20},
		{"invalid", "?limit=abc", "limit", 20, 20},
		{"zero", "?offset=0", "offset", 5, 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test"+tc.queryParam, nil)
			got := QueryInt(req, tc.key, tc.defaultVal)
			if got != tc.expected {
				t.Errorf("QueryInt = %d, want %d", got, tc.expected)
			}
		})
	}
}

func TestDecodeJSON(t *testing.T) {
	body := strings.NewReader(`{"name":"test"}`)
	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", "application/json")

	var result struct {
		Name string `json:"name"`
	}
	if err := DecodeJSON(req, &result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "test" {
		t.Errorf("expected 'test', got %q", result.Name)
	}
}
