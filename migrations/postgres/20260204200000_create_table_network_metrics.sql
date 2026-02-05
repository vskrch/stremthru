-- Network metrics table for storing status checks and speed test results
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS network_metrics (
    id SERIAL PRIMARY KEY,
    type TEXT NOT NULL,         -- 'status' or 'speedtest'
    data JSONB NOT NULL,        -- JSON payload
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_network_metrics_type_created_at ON network_metrics(type, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS network_metrics;
-- +goose StatementEnd
