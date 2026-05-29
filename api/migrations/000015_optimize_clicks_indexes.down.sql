DROP INDEX IF EXISTS idx_clicks_link_id_created_at;

CREATE INDEX idx_clicks_created_at     ON clicks(created_at);
CREATE INDEX idx_clicks_device_type    ON clicks(device_type);
CREATE INDEX idx_clicks_referrer_source ON clicks(referrer_source);
