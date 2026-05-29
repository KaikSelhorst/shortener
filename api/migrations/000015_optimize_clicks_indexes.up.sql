-- Composite index covering the exact filter used by all analytics queries:
-- link_id = $1 AND created_at BETWEEN $2 AND $3
-- This lets Postgres seek directly to the right link and scan only the
-- requested time window instead of loading the full link history.
CREATE INDEX idx_clicks_link_id_created_at ON clicks(link_id, created_at);

-- The separate single-column indexes become redundant for analytics.
-- idx_clicks_link_id is kept because the redirect path writes one click at a
-- time and the FK constraint uses it; idx_clicks_created_at is not needed.
DROP INDEX IF EXISTS idx_clicks_created_at;

-- Low-cardinality indexes on device_type (5 values) and referrer_source
-- (10 values) are never chosen by the planner for GROUP BY aggregations —
-- Postgres prefers a sequential scan filtered on the composite index above.
-- Remove them to save write overhead on every BatchInsert.
DROP INDEX IF EXISTS idx_clicks_device_type;
DROP INDEX IF EXISTS idx_clicks_referrer_source;
