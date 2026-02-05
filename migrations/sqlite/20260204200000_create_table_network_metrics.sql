-- Network metrics table for storing status checks and speed test results
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS network_metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL,         -- 'status' or 'speedtest'
    data TEXT NOT NULL,         -- JSON payload
    created_at DATETIME NOT NULL
);

CREATE INDEX idx_network_metrics_type_created_at ON network_metrics(type, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS network_metrics;
-- +goose StatementEnd
