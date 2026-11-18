package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRoutes checks status codes and bodies of every route.
func TestRoutes(t *testing.T) {
	router := newRouter(slog.New(slog.NewTextHandler(io.Discard, nil)))

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   map[string]string
	}{
		{name: "greeting", path: "/hello/Kiet", wantStatus: http.StatusOK, wantBody: map[string]string{"message": "Hello, Kiet!"}},
		{name: "blank name", path: "/hello/%20", wantStatus: http.StatusBadRequest, wantBody: map[string]string{"error": "name must not be empty"}},
		{name: "unknown route", path: "/nope", wantStatus: http.StatusNotFound},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))

			if rec.Code != tc.wantStatus {
				t.Fatalf("GET %s status = %d, want %d", tc.path, rec.Code, tc.wantStatus)
			}
			if tc.wantBody == nil {
				return
			}
			var got map[string]string
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("GET %s: decode body: %v", tc.path, err)
			}
			for key, want := range tc.wantBody {
				if got[key] != want {
					t.Fatalf("GET %s body[%q] = %q, want %q", tc.path, key, got[key], want)
				}
			}
		})
	}
}

// TestHealthz checks the liveness endpoint.
func TestHealthz(t *testing.T) {
	router := newRouter(slog.New(slog.NewTextHandler(io.Discard, nil)))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusOK || rec.Body.String() != "ok" {
		t.Fatalf("GET /healthz = %d %q, want 200 \"ok\"", rec.Code, rec.Body.String())
	}
}
