package netmon

import (
	"encoding/json"
	"time"

	"github.com/MunifTanjim/stremthru/internal/db"
)

// MetricType identifies the type of stored metric
type MetricType string

const (
	MetricTypeStatus    MetricType = "status"
	MetricTypeSpeedTest MetricType = "speedtest"
)

// StoredMetric represents a persisted metric record
type StoredMetric struct {
	ID        int64      `json:"id"`
	Type      MetricType `json:"type"`
	Data      string     `json:"data"`
	CreatedAt time.Time  `json:"created_at"`
}

// SaveNetworkStatus persists a network status check
func SaveNetworkStatus(status *NetworkStatus) error {
	data, err := json.Marshal(status)
	if err != nil {
		return err
	}
	_, err = db.Exec(
		"INSERT INTO network_metrics (type, data, created_at) VALUES (?, ?, ?)",
		MetricTypeStatus, string(data), status.CheckedAt,
	)
	return err
}

// SaveSpeedTestResult persists a speed test result
func SaveSpeedTestResult(result *SpeedTestResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	_, err = db.Exec(
		"INSERT INTO network_metrics (type, data, created_at) VALUES (?, ?, ?)",
		MetricTypeSpeedTest, string(data), result.TestedAt,
	)
	return err
}

// GetRecentStatuses returns status checks from the last duration
func GetRecentStatuses(since time.Duration) ([]*NetworkStatus, error) {
	cutoff := time.Now().Add(-since)
	rows, err := db.Query(
		"SELECT data FROM network_metrics WHERE type = ? AND created_at > ? ORDER BY created_at DESC LIMIT 100",
		MetricTypeStatus, cutoff,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []*NetworkStatus
	for rows.Next() {
		var data string
		if err := rows.Scan(&data); err != nil {
			continue
		}
		var status NetworkStatus
		if err := json.Unmarshal([]byte(data), &status); err != nil {
			continue
		}
		statuses = append(statuses, &status)
	}
	return statuses, rows.Err()
}

// GetRecentSpeedTests returns speed tests from the last duration
func GetRecentSpeedTests(since time.Duration) ([]*SpeedTestResult, error) {
	cutoff := time.Now().Add(-since)
	rows, err := db.Query(
		"SELECT data FROM network_metrics WHERE type = ? AND created_at > ? ORDER BY created_at DESC LIMIT 50",
		MetricTypeSpeedTest, cutoff,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*SpeedTestResult
	for rows.Next() {
		var data string
		if err := rows.Scan(&data); err != nil {
			continue
		}
		var result SpeedTestResult
		if err := json.Unmarshal([]byte(data), &result); err != nil {
			continue
		}
		results = append(results, &result)
	}
	return results, rows.Err()
}

// CleanupOldMetrics removes metrics older than the specified duration
func CleanupOldMetrics(olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	_, err := db.Exec("DELETE FROM network_metrics WHERE created_at < ?", cutoff)
	return err
}
