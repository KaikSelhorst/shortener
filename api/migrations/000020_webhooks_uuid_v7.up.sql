-- Replace BIGINT identity PK on webhooks with UUID,
-- and update the webhook_id FK on webhook_deliveries accordingly.

-- 1. Add temporary UUID column to webhooks for existing rows
ALTER TABLE webhooks ADD COLUMN new_id UUID NOT NULL DEFAULT gen_random_uuid();

-- 2. Add temporary UUID column to webhook_deliveries
ALTER TABLE webhook_deliveries ADD COLUMN new_webhook_id UUID;

-- 3. Map old integer IDs to the new UUIDs
UPDATE webhook_deliveries wd
SET new_webhook_id = w.new_id
FROM webhooks w
WHERE wd.webhook_id = w.id;

-- 4. Drop old FK and PK constraints
ALTER TABLE webhook_deliveries DROP CONSTRAINT webhook_deliveries_webhook_id_fkey;
ALTER TABLE webhooks DROP CONSTRAINT webhooks_pkey;

-- 5. Swap columns on webhooks
ALTER TABLE webhooks DROP COLUMN id;
ALTER TABLE webhooks RENAME COLUMN new_id TO id;
ALTER TABLE webhooks ADD PRIMARY KEY (id);

-- 6. Swap columns on webhook_deliveries
ALTER TABLE webhook_deliveries DROP COLUMN webhook_id;
ALTER TABLE webhook_deliveries RENAME COLUMN new_webhook_id TO webhook_id;
ALTER TABLE webhook_deliveries ALTER COLUMN webhook_id SET NOT NULL;
ALTER TABLE webhook_deliveries ADD CONSTRAINT webhook_deliveries_webhook_id_fkey
    FOREIGN KEY (webhook_id) REFERENCES webhooks(id) ON DELETE CASCADE;

-- 7. Recreate index
DROP INDEX IF EXISTS idx_webhooks_project_id;
CREATE INDEX idx_webhooks_project_id ON webhooks(project_id);
