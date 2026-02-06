CREATE TABLE IF NOT EXISTS lesson_sections (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  lesson_id     UUID      NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
  title         VARCHAR(255) NOT NULL,
  position      INT         NOT NULL DEFAULT 0,
  is_published  BOOLEAN     NOT NULL DEFAULT TRUE,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT uidx_lesson_section_pos UNIQUE (lesson_id, position)
);

CREATE INDEX IF NOT EXISTS idx_lesson_sections_lesson_id ON lesson_sections (lesson_id);