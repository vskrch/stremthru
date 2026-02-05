package netmon

import (
	"testing"
)

func TestNewMonitor(t *testing.T) {
	m := NewMonitor()
	if m == nil {
		t.Fatal("NewMonitor returned nil")
	}
	if m.startTime.IsZero() {
		t.Error("startTime is zero")
	}
}

func TestDefault(t *testing.T) {
	m := Default()
	if m == nil {
		t.Fatal("Default returned nil")
	}
}

func TestBandwidthMetrics(t *testing.T) {
	m := NewMonitor()

	m.AddBytesIn(100)
	m.AddBytesOut(200)
	m.AddSpeedTestBytes(300)

	metrics := m.GetBandwidth()

	if metrics.DashboardBytesIn != 100 {
		t.Errorf("expected 100 bytes in, got %d", metrics.DashboardBytesIn)
	}
	if metrics.DashboardBytesOut != 200 {
		t.Errorf("expected 200 bytes out, got %d", metrics.DashboardBytesOut)
	}
	if metrics.SpeedTestBytes != 300 {
		t.Errorf("expected 300 speed test bytes, got %d", metrics.SpeedTestBytes)
	}
	if metrics.SinceRestartAt != m.startTime {
		t.Error("start time mismatch")
	}
}

func TestGetCachedStatus(t *testing.T) {
	m := NewMonitor()
	status := &NetworkStatus{MachineIP: "1.2.3.4"}

	m.mu.Lock()
	m.lastStatus = status
	m.mu.Unlock()

	cached := m.GetCachedStatus()
	if cached != status {
		t.Error("cached status mismatch")
	}
}
