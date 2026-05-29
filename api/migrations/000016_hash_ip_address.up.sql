-- Replace the raw ip_address (INET, personally identifiable) with a
-- pseudonymous HMAC-SHA256 hash computed at record time using a server-side
-- secret (IP_HASH_SECRET). The original IP address is irrecoverable without
-- that secret, satisfying GDPR/LGPD pseudonymisation requirements while still
-- allowing COUNT(DISTINCT ip_hash) for unique-visitor analytics.
--
-- Existing rows have their raw IP deleted — they cannot be retroactively
-- hashed because the secret was not available at insert time.
ALTER TABLE clicks
  ADD COLUMN ip_hash TEXT;

ALTER TABLE clicks
  DROP COLUMN ip_address;
