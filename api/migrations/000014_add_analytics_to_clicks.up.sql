ALTER TABLE clicks
  ADD COLUMN device_type     TEXT NOT NULL DEFAULT 'unknown',
  ADD COLUMN referrer_source TEXT NOT NULL DEFAULT 'direct';

CREATE INDEX idx_clicks_device_type     ON clicks(device_type);
CREATE INDEX idx_clicks_referrer_source ON clicks(referrer_source);
