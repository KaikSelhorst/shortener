ALTER TABLE clicks
  ADD COLUMN ip_address INET;

ALTER TABLE clicks
  DROP COLUMN ip_hash;
