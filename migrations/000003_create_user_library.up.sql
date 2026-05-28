CREATE TABLE user_library (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  media_id    UUID NOT NULL REFERENCES media_cache(id),
  status      TEXT NOT NULL CHECK (status IN ('planning','in_progress','completed','dropped','on_hold')),
  user_rating SMALLINT CHECK (user_rating BETWEEN 1 AND 10),
  notes       TEXT,
  started_at  DATE,
  finished_at DATE,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at  TIMESTAMPTZ,
  CONSTRAINT user_library_unique UNIQUE (user_id, media_id)
);
CREATE INDEX idx_library_user        ON user_library(user_id)         WHERE deleted_at IS NULL;
CREATE INDEX idx_library_user_status ON user_library(user_id, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_library_media       ON user_library(media_id);
