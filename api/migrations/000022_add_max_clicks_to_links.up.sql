ALTER TABLE links ADD COLUMN max_clicks BIGINT CHECK (max_clicks > 0);
