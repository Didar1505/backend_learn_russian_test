CREATE TABLE IF NOT EXISTS lesson_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  lesson_section_id  UUID      NOT NULL REFERENCES lesson_sections(id) ON DELETE CASCADE,
  kind              VARCHAR(50)  NOT NULL,
  position          INT          NOT NULL DEFAULT 0,
  payload           JSONB        NOT NULL,
  meta              JSONB        NOT NULL DEFAULT '{}'::jsonb,
  is_published      BOOLEAN      NOT NULL DEFAULT TRUE,
  created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  updated_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  CONSTRAINT uidx_section_item_pos UNIQUE (lesson_section_id, position),
  CONSTRAINT chk_lesson_items_kind_not_empty CHECK (length(trim(kind)) > 0)
);

CREATE INDEX IF NOT EXISTS idx_lesson_items_section_id ON lesson_items (lesson_section_id);
CREATE INDEX IF NOT EXISTS idx_lesson_items_kind       ON lesson_items (kind);

-- Useful for JSON queries later (optional but nice)
CREATE INDEX IF NOT EXISTS idx_lesson_items_payload_gin ON lesson_items USING GIN (payload);