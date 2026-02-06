CREATE TABLE IF NOT EXISTS courses (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  slug          VARCHAR(120) NOT NULL UNIQUE,
  title         VARCHAR(255) NOT NULL,
  description   TEXT,
  level         VARCHAR(10)  NOT NULL,
  language_from VARCHAR(10)  NOT NULL,
  language_to   VARCHAR(10)  NOT NULL,
  is_published  BOOLEAN      NOT NULL DEFAULT FALSE,
  created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_courses_is_published ON courses (is_published);
CREATE INDEX IF NOT EXISTS idx_courses_level        ON courses (level);