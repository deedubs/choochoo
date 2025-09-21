-- Create webhook_events table to store GitHub webhook payloads
CREATE TABLE webhook_events (
    id SERIAL PRIMARY KEY,
    delivery_id VARCHAR(255) NOT NULL UNIQUE,
    event_type VARCHAR(50) NOT NULL,
    repository_name VARCHAR(255),
    sender_login VARCHAR(255),
    action VARCHAR(100),
    payload JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Add indexes for common queries
CREATE INDEX idx_webhook_events_event_type ON webhook_events (event_type);
CREATE INDEX idx_webhook_events_repository ON webhook_events (repository_name);
CREATE INDEX idx_webhook_events_sender ON webhook_events (sender_login);
CREATE INDEX idx_webhook_events_created_at ON webhook_events (created_at);

-- Add a comment to the table
COMMENT ON TABLE webhook_events IS 'Stores GitHub webhook events for push, issue_comment, and pull_request events';