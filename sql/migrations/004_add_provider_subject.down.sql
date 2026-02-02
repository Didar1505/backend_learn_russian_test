DROP INDEX IF EXISTS uniq_users_provider_subject;
ALTER TABLE users DROP COLUMN IF EXISTS provider_subject;