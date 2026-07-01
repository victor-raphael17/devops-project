package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestProjetoKorpEndpoint(t *testing.T) {
	srv := httptest.NewServer(newRouter())
	defer srv.Close()

	res, err := http.Get(srv.URL + "/projeto-korp")
	if err != nil {
		t.Fatalf("GET /projeto-korp: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.StatusCode, http.StatusOK)
	}
	if ct := res.Header.Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %q, want application/json; charset=utf-8", ct)
	}

	var body response
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decoding body: %v", err)
	}

	if body.Name != projectName {
		t.Errorf("nome = %q, want %q", body.Name, projectName)
	}
	if !strings.HasSuffix(body.Time, "Z") {
		t.Errorf("horario = %q, want UTC (Z suffix)", body.Time)
	}

	parsed, err := time.Parse(time.RFC3339, body.Time)
	if err != nil {
		t.Fatalf("horario %q is not RFC 3339: %v", body.Time, err)
	}
	if since := time.Since(parsed); since < 0 || since > time.Minute {
		t.Errorf("horario %q is not current (off by %v)", body.Time, since)
	}
}

func TestRouting(t *testing.T) {
	srv := httptest.NewServer(newRouter())
	defer srv.Close()

	tests := []struct {
		method string
		path   string
		want   int
	}{
		{http.MethodGet, "/healthz", http.StatusOK},
		{http.MethodGet, "/metrics", http.StatusOK},
		{http.MethodPost, "/projeto-korp", http.StatusMethodNotAllowed},
		{http.MethodGet, "/nao-existe", http.StatusNotFound},
	}

	for _, tc := range tests {
		req, err := http.NewRequest(tc.method, srv.URL+tc.path, nil)
		if err != nil {
			t.Fatalf("%s %s: building request: %v", tc.method, tc.path, err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("%s %s: %v", tc.method, tc.path, err)
		}
		res.Body.Close()

		if res.StatusCode != tc.want {
			t.Errorf("%s %s: status = %d, want %d", tc.method, tc.path, res.StatusCode, tc.want)
		}
	}
}

func TestPort(t *testing.T) {
	t.Setenv("PORT", "")
	if got := port(); got != defaultPort {
		t.Errorf("port() with PORT unset = %q, want %q", got, defaultPort)
	}

	t.Setenv("PORT", "9999")
	if got := port(); got != "9999" {
		t.Errorf("port() with PORT=9999 = %q, want %q", got, "9999")
	}
}
