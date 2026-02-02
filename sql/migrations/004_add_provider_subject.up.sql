ALTER TABLE users ADD COLUMN provider_subject TEXT;

CREATE UNIQUE INDEX uniq_users_provider_subject
ON users(auth_provider, provider_subject)
WHERE provider_subject IS NOT NULL;
