package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/CristianSsousa/saga-api/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type libraryRepo struct{ db *pgxpool.Pool }

func NewLibraryRepo(db *pgxpool.Pool) domain.LibraryRepository {
	return &libraryRepo{db: db}
}

func (r *libraryRepo) Add(ctx context.Context, userID, mediaID string, status domain.Status) (*domain.LibraryItem, error) {
	var id string
	err := r.db.QueryRow(ctx,
		`INSERT INTO user_library (user_id, media_id, status)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (user_id, media_id) DO UPDATE
		   SET status     = EXCLUDED.status,
		       deleted_at = NULL,
		       updated_at = now()
		 RETURNING id`,
		userID, mediaID, string(status),
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id, userID)
}

func (r *libraryRepo) Update(ctx context.Context, id, userID string, status *domain.Status, rating *int, notes *string) (*domain.LibraryItem, error) {
	setClauses := []string{}
	args := []any{}
	argIdx := 1

	if status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, string(*status))
		argIdx++
	}
	if rating != nil {
		setClauses = append(setClauses, fmt.Sprintf("user_rating = $%d", argIdx))
		args = append(args, *rating)
		argIdx++
	}
	if notes != nil {
		setClauses = append(setClauses, fmt.Sprintf("notes = $%d", argIdx))
		args = append(args, *notes)
		argIdx++
	}

	if len(setClauses) == 0 {
		return r.FindByID(ctx, id, userID)
	}

	args = append(args, id, userID)
	query := fmt.Sprintf(
		`UPDATE user_library SET %s WHERE id = $%d AND user_id = $%d AND deleted_at IS NULL`,
		strings.Join(setClauses, ", "),
		argIdx,
		argIdx+1,
	)

	tag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, domain.ErrLibraryItemNotFound
	}
	return r.FindByID(ctx, id, userID)
}

func (r *libraryRepo) SoftDelete(ctx context.Context, id, userID string) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE user_library SET deleted_at = now() WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`,
		id, userID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrLibraryItemNotFound
	}
	return nil
}

func (r *libraryRepo) List(ctx context.Context, userID string, filter domain.ListFilter) ([]domain.LibraryItem, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}

	args := []any{userID}
	argIdx := 2
	conditions := []string{"ul.user_id = $1", "ul.deleted_at IS NULL"}

	if filter.Cursor != "" {
		conditions = append(conditions, fmt.Sprintf("ul.id > $%d", argIdx))
		args = append(args, filter.Cursor)
		argIdx++
	}
	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("ul.status = $%d", argIdx))
		args = append(args, string(*filter.Status))
		argIdx++
	}
	if filter.MediaType != nil {
		conditions = append(conditions, fmt.Sprintf("mc.media_type = $%d", argIdx))
		args = append(args, string(*filter.MediaType))
		argIdx++
	}

	args = append(args, limit)

	query := fmt.Sprintf(`
		SELECT
			ul.id, ul.user_id, ul.status, ul.user_rating, ul.notes,
			ul.started_at, ul.finished_at, ul.created_at, ul.updated_at,
			mc.external_id, mc.source, mc.media_type, mc.title, mc.cover_url, mc.metadata
		FROM user_library ul
		JOIN media_cache mc ON mc.id = ul.media_id
		WHERE %s
		ORDER BY ul.id ASC
		LIMIT $%d`,
		strings.Join(conditions, " AND "),
		argIdx,
	)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanLibraryRows(rows)
}

func (r *libraryRepo) FindByID(ctx context.Context, id, userID string) (*domain.LibraryItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			ul.id, ul.user_id, ul.status, ul.user_rating, ul.notes,
			ul.started_at, ul.finished_at, ul.created_at, ul.updated_at,
			mc.external_id, mc.source, mc.media_type, mc.title, mc.cover_url, mc.metadata
		FROM user_library ul
		JOIN media_cache mc ON mc.id = ul.media_id
		WHERE ul.id = $1 AND ul.user_id = $2 AND ul.deleted_at IS NULL`,
		id, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items, err := scanLibraryRows(rows)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, domain.ErrLibraryItemNotFound
	}
	return &items[0], nil
}

func scanLibraryRows(rows pgx.Rows) ([]domain.LibraryItem, error) {
	var items []domain.LibraryItem
	for rows.Next() {
		var item domain.LibraryItem
		var statusStr, mediaTypeStr string
		var metaRaw []byte
		var startedAt, finishedAt *string

		err := rows.Scan(
			&item.ID, &item.UserID, &statusStr, &item.UserRating, &item.Notes,
			&startedAt, &finishedAt, &item.CreatedAt, &item.UpdatedAt,
			&item.Media.ExternalID, &item.Media.Source, &mediaTypeStr,
			&item.Media.Title, &item.Media.CoverURL, &metaRaw,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, domain.ErrLibraryItemNotFound
			}
			return nil, err
		}
		item.Status = domain.Status(statusStr)
		item.Media.Type = domain.MediaType(mediaTypeStr)
		item.StartedAt = startedAt
		item.FinishedAt = finishedAt
		if metaRaw != nil {
			_ = json.Unmarshal(metaRaw, &item.Media.Metadata)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
