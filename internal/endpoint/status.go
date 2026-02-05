package endpoint

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/MunifTanjim/stremthru/internal/config"
	"github.com/MunifTanjim/stremthru/internal/netmon"
)

var statusMonitor = netmon.Default()

// adminAuthMiddleware protects status endpoints with admin credentials
func adminAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="StremThru Status"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		expectedPass := config.AdminPassword.GetPassword(user)
		if expectedPass == "" || expectedPass != pass {
			w.Header().Set("WWW-Authenticate", `Basic realm="StremThru Status"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// handleStatusPage serves the status dashboard HTML
func handleStatusPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.New("status").Parse(statusHTML)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, map[string]string{
		"Version": config.Version,
	})
}

// handleStatusCurrent returns current network status
func handleStatusCurrent(w http.ResponseWriter, r *http.Request) {
	status, err := statusMonitor.GetStatus()
	if err != nil {
		sendStatusError(w, err)
		return
	}

	// Store the status (ignore errors - storage is optional)
	netmon.SaveNetworkStatus(status)

	sendStatusJSON(w, status)
}

// handleStatusSpeedTest triggers a speed test
func handleStatusSpeedTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get RD API key from config if available
	rdAPIKey := config.StoreAuthToken.GetToken("*", "realdebrid")

	result, err := statusMonitor.RunSpeedTest(rdAPIKey)
	if err != nil {
		sendStatusError(w, err)
		return
	}

	// Store the result (ignore errors - storage is optional)
	netmon.SaveSpeedTestResult(result)

	sendStatusJSON(w, result)
}

// handleStatusHistory returns historical metrics
func handleStatusHistory(w http.ResponseWriter, r *http.Request) {
	// Get data from last 24 hours
	duration := 24 * time.Hour
	if r.URL.Query().Get("range") == "7d" {
		duration = 7 * 24 * time.Hour
	}

	statuses, _ := netmon.GetRecentStatuses(duration)
	speedTests, _ := netmon.GetRecentSpeedTests(duration)

	// Handle nil slices for JSON
	if statuses == nil {
		statuses = []*netmon.NetworkStatus{}
	}
	if speedTests == nil {
		speedTests = []*netmon.SpeedTestResult{}
	}

	sendStatusJSON(w, map[string]interface{}{
		"statuses":    statuses,
		"speed_tests": speedTests,
	})
}

// handleStatusBandwidth returns bandwidth metrics
func handleStatusBandwidth(w http.ResponseWriter, r *http.Request) {
	bandwidth := statusMonitor.GetBandwidth()
	sendStatusJSON(w, bandwidth)
}

// handleStatusPing returns a simple ping for latency measurement
func handleStatusPing(w http.ResponseWriter, r *http.Request) {
	sendStatusJSON(w, map[string]interface{}{
		"pong":      true,
		"timestamp": time.Now().UnixMilli(),
		"server_ip": config.IP.GetMachineIP(),
	})
}

func sendStatusJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func sendStatusError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

// AddStatusEndpoints registers all status dashboard routes
func AddStatusEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("/status", adminAuthMiddleware(handleStatusPage))
	mux.HandleFunc("/status/api/current", adminAuthMiddleware(handleStatusCurrent))
	mux.HandleFunc("/status/api/speedtest", adminAuthMiddleware(handleStatusSpeedTest))
	mux.HandleFunc("/status/api/history", adminAuthMiddleware(handleStatusHistory))
	mux.HandleFunc("/status/api/bandwidth", adminAuthMiddleware(handleStatusBandwidth))
	mux.HandleFunc("/status/api/ping", adminAuthMiddleware(handleStatusPing))
}
