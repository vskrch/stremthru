package netmon

import (
	"io"
	"net/http"
	"time"

	"github.com/MunifTanjim/stremthru/internal/config"
)

// SpeedTestSize is the approximate size to download for speed tests
const SpeedTestSize = 10 * 1024 * 1024 // 10MB

// RunSpeedTest performs a multi-segment speed test
func (m *Monitor) RunSpeedTest(rdAPIKey string) (*SpeedTestResult, error) {
	result := &SpeedTestResult{
		TestedAt: time.Now(),
	}

	totalLatency := 0

	// Test 1: Server → WARP (Cloudflare speed test)
	warpResult, err := m.testServerToWarp()
	if err != nil {
		result.Error = "warp: " + err.Error()
	} else {
		result.ServerToWarp = warpResult
		totalLatency += warpResult.LatencyMs
	}

	// Test 2: WARP → RealDebrid (download test file)
	if rdAPIKey != "" {
		rdResult, err := m.testWarpToRealDebrid(rdAPIKey)
		if err != nil {
			if result.Error == "" {
				result.Error = "realdebrid: " + err.Error()
			}
		} else {
			result.WarpToRealDebrid = rdResult
			totalLatency += rdResult.LatencyMs
		}
	}

	result.TotalLatencyMs = totalLatency

	// Cache the result
	m.mu.Lock()
	m.lastSpeedTest = result
	m.mu.Unlock()

	return result, nil
}

// testServerToWarp tests speed from server to Cloudflare via WARP tunnel
func (m *Monitor) testServerToWarp() (*SegmentResult, error) {
	client := config.GetHTTPClient(config.TUNNEL_TYPE_FORCED)
	client.Timeout = 30 * time.Second

	// Use Cloudflare's speed test endpoint
	url := "https://speed.cloudflare.com/__down?bytes=5000000" // 5MB

	start := time.Now()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// First byte latency
	latency := int(time.Since(start).Milliseconds())

	// Read all data and measure
	downloadStart := time.Now()
	written, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		return nil, err
	}
	downloadDuration := time.Since(downloadStart).Seconds()

	// Track bandwidth
	m.AddSpeedTestBytes(written)

	// Calculate speed in Mbps
	speedMbps := float64(written*8) / downloadDuration / 1_000_000

	// Get IPs
	machineIP := config.IP.GetMachineIP()
	tunnelIP, _ := config.IP.GetTunnelIP()

	return &SegmentResult{
		SpeedMbps:        speedMbps,
		LatencyMs:        latency,
		BytesTransferred: written,
		SourceIP:         machineIP,
		DestIP:           tunnelIP,
	}, nil
}

// testWarpToRealDebrid tests speed from WARP to RealDebrid
func (m *Monitor) testWarpToRealDebrid(apiKey string) (*SegmentResult, error) {
	client := config.GetHTTPClient(config.TUNNEL_TYPE_FORCED)
	client.Timeout = 60 * time.Second

	// Get a small test download from RD - use their streaming test
	// First, we need to get a download URL from RD's test endpoint
	testURL := "https://real-debrid.com/speedtest"

	start := time.Now()

	req, err := http.NewRequest(http.MethodGet, testURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// First byte latency
	latency := int(time.Since(start).Milliseconds())

	// Read response and measure
	downloadStart := time.Now()
	written, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		return nil, err
	}
	downloadDuration := time.Since(downloadStart).Seconds()

	// Track bandwidth
	m.AddSpeedTestBytes(written)

	// Calculate speed (handle very small downloads)
	var speedMbps float64
	if downloadDuration > 0 && written > 0 {
		speedMbps = float64(written*8) / downloadDuration / 1_000_000
	}

	// Get tunnel IP as source (what RD sees)
	tunnelIP, _ := config.IP.GetTunnelIP()

	return &SegmentResult{
		SpeedMbps:        speedMbps,
		LatencyMs:        latency,
		BytesTransferred: written,
		SourceIP:         tunnelIP,
		DestIP:           "real-debrid.com",
	}, nil
}

// GetLastSpeedTest returns the most recent speed test result
func (m *Monitor) GetLastSpeedTest() *SpeedTestResult {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastSpeedTest
}
