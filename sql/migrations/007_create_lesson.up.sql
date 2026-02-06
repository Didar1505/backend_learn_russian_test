CREATE TABLE IF NOT EXISTS lessons (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  course_id          UUID      NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  module_id          UUID      NULL REFERENCES modules(id) ON DELETE SET NULL,
  title              VARCHAR(255) NOT NULL,
  summary            TEXT,
  position           INT         NOT NULL DEFAULT 0,
  estimated_minutes  INT         NOT NULL DEFAULT 0,
  is_published       BOOLEAN     NOT NULL DEFAULT FALSE,
  created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_lessons_course_id ON lessons (course_id);
CREATE INDEX IF NOT EXISTS idx_lessons_module_id ON lessons (module_id);
CREATE INDEX IF NOT EXISTS idx_lessons_order_course ON lessons (course_id, position);
CREATE INDEX IF NOT EXISTS idx_lessons_order_module ON lessons (module_id, position);
