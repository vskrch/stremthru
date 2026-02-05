package endpoint

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MunifTanjim/stremthru/internal/config"
	"github.com/MunifTanjim/stremthru/internal/netmon"
)

func TestAdminAuthMiddleware(t *testing.T) {
	// Setup config for test
	config.AdminPassword = config.UserPasswordMap{
		"testuser": "testpass",
	}

	handler := adminAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name       string
		user       string
		pass       string
		wantStatus int
	}{
		{"NoAuth", "", "", http.StatusUnauthorized},
		{"WrongUser", "wrong", "testpass", http.StatusUnauthorized},
		{"WrongPass", "testuser", "wrong", http.StatusUnauthorized},
		{"CorrectAuth", "testuser", "testpass", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.user != "" {
				req.SetBasicAuth(tt.user, tt.pass)
			}
			w := httptest.NewRecorder()
			handler(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestHandleStatusPage(t *testing.T) {
	req := httptest.NewRequest("GET", "/status", nil)
	w := httptest.NewRecorder()

	handleStatusPage(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}
	if w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Error("content type is not text/html")
	}
}

// Mock netmon functions/structs would be needed for full coverage of API endpoints
// but since we're using global/package state, it's harder to mock without refactoring.
// For now, testing the Ping endpoint which is self-contained.

func TestHandleStatusPing(t *testing.T) {
	req := httptest.NewRequest("GET", "/status/api/ping", nil)
	w := httptest.NewRecorder()

	handleStatusPing(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp["pong"] != true {
		t.Error("pong is not true")
	}
}

func TestHandleStatusBandwidth(t *testing.T) {
	req := httptest.NewRequest("GET", "/status/api/bandwidth", nil)
	w := httptest.NewRecorder()

	handleStatusBandwidth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}

	var resp netmon.BandwidthMetrics
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
}
