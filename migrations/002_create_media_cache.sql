CREATE TABLE media_cache (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  external_id TEXT NOT NULL,
  source      TEXT NOT NULL CHECK (source IN ('tmdb','rawg','openlibrary','anilist')),
  media_type  TEXT NOT NULL CHECK (media_type IN ('movie','tv','game','book','manga')),
  title       TEXT NOT NULL,
  cover_url   TEXT,
  metadata    JSONB NOT NULL DEFAULT '{}',
  cached_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at  TIMESTAMPTZ,
  CONSTRAINT media_cache_source_external UNIQUE (source, external_id)
);
CREATE INDEX idx_media_cache_type   ON media_cache(media_type);
CREATE INDEX idx_media_cache_source ON media_cache(source, external_id);
CREATE INDEX idx_media_cache_expiry ON media_cache(expires_at) WHERE expires_at IS NOT NULL;
