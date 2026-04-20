DROP INDEX IF EXISTS idx_links_short_code;

ALTER TABLE links DROP COLUMN short_code;
