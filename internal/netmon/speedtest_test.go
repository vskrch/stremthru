package netmon

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRunSpeedTest(t *testing.T) {
	// Start a local test server to simulate download
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve 1MB of 0s
		data := make([]byte, 1024*1024)
		w.Header().Set("Content-Length", "1048576")
		w.Write(data)
	}))
	defer ts.Close()

	// We can't easily mock the external calls in RunSpeedTest without refactoring,
	// but we can test the helper functions if we export them or test internals.
	// For now, let's test the throughput calculation logic by using downloadFile directly if possible
	// or create a focused test that doesn't hit external APIs.

	// Since downloadFile is private/internal in speedtest.go, we test what we can access.
	// Let's create a test that calls downloadFile (whitebox testing since we are in same package)

	start := time.Now()
	size, err := downloadFile(http.DefaultClient, ts.URL)
	if err != nil {
		t.Fatalf("downloadFile failed: %v", err)
	}
	duration := time.Since(start)

	if size != 1048576 {
		t.Errorf("expected 1MB size, got %d", size)
	}

	// Calculate speed
	speedMbps := float64(size*8) / (float64(duration.Milliseconds()) / 1000.0) / 1000000.0
	log.Printf("Test download speed: %.2f Mbps", speedMbps)
}

func TestDownloadFile_Error(t *testing.T) {
	// Test error handling
	_, err := downloadFile(http.DefaultClient, "http://invalid-url-that-does-not-exist.local")
	if err == nil {
		t.Error("expected error for invalid URL")
	}
}
