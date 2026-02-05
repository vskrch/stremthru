package netmon

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/MunifTanjim/stremthru/internal/config"
)

// NetworkStatus contains current IP and connectivity information
type NetworkStatus struct {
	MachineIP         string    `json:"machine_ip"`
	TunnelIP          string    `json:"tunnel_ip"`
	TunnelASN         string    `json:"tunnel_asn,omitempty"`
	TunnelOrg         string    `json:"tunnel_org,omitempty"`
	WarpActive        bool      `json:"warp_active"`
	RealDebridOK      bool      `json:"realdebrid_ok"`
	RealDebridLatency int       `json:"realdebrid_latency_ms"`
	RealDebridSeenIP  string    `json:"realdebrid_seen_ip,omitempty"` // IP that RD sees
	LastError         string    `json:"last_error,omitempty"`
	CheckedAt         time.Time `json:"checked_at"`
}

// SegmentResult represents speed test results for one network segment
type SegmentResult struct {
	SpeedMbps        float64 `json:"speed_mbps"`
	LatencyMs        int     `json:"latency_ms"`
	BytesTransferred int64   `json:"bytes_transferred"`
	SourceIP         string  `json:"source_ip"`
	DestIP           string  `json:"dest_ip,omitempty"`
}

// SpeedTestResult contains multi-hop speed test results
type SpeedTestResult struct {
	ServerToWarp     *SegmentResult `json:"server_to_warp,omitempty"`
	WarpToRealDebrid *SegmentResult `json:"warp_to_realdebrid,omitempty"`
	TotalLatencyMs   int            `json:"total_latency_ms"`
	TestedAt         time.Time      `json:"tested_at"`
	Error            string         `json:"error,omitempty"`
}

// BandwidthMetrics tracks data transfer counters
type BandwidthMetrics struct {
	DashboardBytesIn  int64     `json:"dashboard_bytes_in"`
	DashboardBytesOut int64     `json:"dashboard_bytes_out"`
	SpeedTestBytes    int64     `json:"speedtest_bytes"`
	SinceRestartAt    time.Time `json:"since_restart_at"`
}

// Monitor handles all network monitoring operations
type Monitor struct {
	mu             sync.RWMutex
	lastStatus     *NetworkStatus
	lastSpeedTest  *SpeedTestResult
	bytesIn        atomic.Int64
	bytesOut       atomic.Int64
	speedTestBytes atomic.Int64
	startTime      time.Time
}

// NewMonitor creates a new network monitor
func NewMonitor() *Monitor {
	return &Monitor{
		startTime: time.Now(),
	}
}

var defaultMonitor = NewMonitor()

// Default returns the default monitor instance
func Default() *Monitor {
	return defaultMonitor
}

// GetStatus performs network status check
func (m *Monitor) GetStatus() (*NetworkStatus, error) {
	status := &NetworkStatus{
		CheckedAt: time.Now(),
	}

	// Get machine IP (direct, no tunnel)
	status.MachineIP = config.IP.GetMachineIP()

	// Get tunnel IP (through WARP)
	tunnelIP, err := config.IP.GetTunnelIP()
	if err != nil {
		status.LastError = "tunnel: " + err.Error()
	} else {
		status.TunnelIP = tunnelIP
		// Check if different from machine IP = WARP is active
		status.WarpActive = tunnelIP != "" && tunnelIP != status.MachineIP
	}

	// Get ASN info for tunnel IP
	if status.TunnelIP != "" {
		asn, org := getASNInfo(status.TunnelIP)
		status.TunnelASN = asn
		status.TunnelOrg = org
	}

	// Check RealDebrid connectivity
	rdOK, rdLatency, rdSeenIP, rdErr := checkRealDebrid()
	status.RealDebridOK = rdOK
	status.RealDebridLatency = rdLatency
	status.RealDebridSeenIP = rdSeenIP
	if rdErr != nil && status.LastError == "" {
		status.LastError = "realdebrid: " + rdErr.Error()
	}

	// Cache the status
	m.mu.Lock()
	m.lastStatus = status
	m.mu.Unlock()

	return status, nil
}

// GetCachedStatus returns the last known status without re-checking
func (m *Monitor) GetCachedStatus() *NetworkStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastStatus
}

// GetBandwidth returns current bandwidth metrics
func (m *Monitor) GetBandwidth() *BandwidthMetrics {
	return &BandwidthMetrics{
		DashboardBytesIn:  m.bytesIn.Load(),
		DashboardBytesOut: m.bytesOut.Load(),
		SpeedTestBytes:    m.speedTestBytes.Load(),
		SinceRestartAt:    m.startTime,
	}
}

// AddBytesIn increments the bytes received counter
func (m *Monitor) AddBytesIn(n int64) {
	m.bytesIn.Add(n)
}

// AddBytesOut increments the bytes sent counter
func (m *Monitor) AddBytesOut(n int64) {
	m.bytesOut.Add(n)
}

// AddSpeedTestBytes increments the speed test bytes counter
func (m *Monitor) AddSpeedTestBytes(n int64) {
	m.speedTestBytes.Add(n)
}

// getASNInfo fetches ASN information for an IP (simplified)
func getASNInfo(ip string) (asn string, org string) {
	// Use ip-api.com for quick ASN lookup (rate limited, but fine for occasional use)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://ip-api.com/json/" + ip + "?fields=as,org")
	if err != nil {
		return "", ""
	}
	defer resp.Body.Close()

	var result struct {
		AS  string `json:"as"`
		Org string `json:"org"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", ""
	}
	return result.AS, result.Org
}

// checkRealDebrid tests connectivity to RealDebrid and gets the IP they see
func checkRealDebrid() (ok bool, latencyMs int, seenIP string, err error) {
	// Use the tunnel client to check what IP external services see
	client := config.GetHTTPClient(config.TUNNEL_TYPE_FORCED)
	client.Timeout = 10 * time.Second

	// First, get the IP that external services see when we use the tunnel
	// This is the IP that RealDebrid will see
	seenIP = getIPViaTunnel(client)

	start := time.Now()

	// Check RD API connectivity
	req, err := http.NewRequest(http.MethodGet, "https://api.real-debrid.com/time", nil)
	if err != nil {
		return false, 0, seenIP, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, 0, seenIP, err
	}
	defer resp.Body.Close()

	latencyMs = int(time.Since(start).Milliseconds())

	if resp.StatusCode != http.StatusOK {
		return false, latencyMs, seenIP, nil
	}

	// Read response (just to complete the request)
	io.Copy(io.Discard, resp.Body)

	return true, latencyMs, seenIP, nil
}

// getIPViaTunnel gets the public IP as seen by external services through the tunnel
func getIPViaTunnel(client *http.Client) string {
	// Use checkip.amazonaws.com which is reliable and fast
	req, err := http.NewRequest(http.MethodGet, "https://checkip.amazonaws.com", nil)
	if err != nil {
		return ""
	}

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(body))
}
