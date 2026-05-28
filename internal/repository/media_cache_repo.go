package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/CristianSsousa/saga-api/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MediaCacheRepo struct{ db *pgxpool.Pool }

func NewMediaCacheRepo(db *pgxpool.Pool) *MediaCacheRepo {
	return &MediaCacheRepo{db: db}
}

// Upsert stores or updates media metadata and returns the internal UUID.
func (r *MediaCacheRepo) Upsert(ctx context.Context, m domain.Media) (string, error) {
	meta, err := json.Marshal(m.Metadata)
	if err != nil {
		return "", err
	}
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	var id string
	err = r.db.QueryRow(ctx, `
		INSERT INTO media_cache (external_id, source, media_type, title, cover_url, metadata, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (source, external_id) DO UPDATE
		  SET title = EXCLUDED.title,
		      cover_url = EXCLUDED.cover_url,
		      metadata = EXCLUDED.metadata,
		      cached_at = now(),
		      expires_at = EXCLUDED.expires_at
		RETURNING id`,
		m.ExternalID, m.Source, string(m.Type), m.Title, m.CoverURL, meta, expiresAt,
	).Scan(&id)
	return id, err
}

// FindByID returns a Media struct from the cache by internal UUID.
func (r *MediaCacheRepo) FindByID(ctx context.Context, id string) (*domain.Media, error) {
	var m domain.Media
	var mt string
	var metaRaw []byte
	err := r.db.QueryRow(ctx,
		`SELECT external_id, source, media_type, title, cover_url, metadata FROM media_cache WHERE id = $1`,
		id,
	).Scan(&m.ExternalID, &m.Source, &mt, &m.Title, &m.CoverURL, &metaRaw)
	if err != nil {
		return nil, err
	}
	m.Type = domain.MediaType(mt)
	_ = json.Unmarshal(metaRaw, &m.Metadata)
	return &m, nil
}
