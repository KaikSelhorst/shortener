ALTER TABLE links ADD COLUMN short_code VARCHAR(20) NOT NULL UNIQUE;

CREATE INDEX idx_links_short_code ON links(short_code);
