package repository

import (
	"context"
	"database/sql"
	"rsslibrary/internal/app/domain"
	"time"
)

type FeedsRepositoryInterface interface {
	InsertFeed(ctx context.Context, name, url string) (domain.Feed, error)
	GetAllFeeds(ctx context.Context) ([]domain.Feed, error)
	GetFeedByName(ctx context.Context, name string) (domain.Feed, error)
	GetOldestFeeds(ctx context.Context, limit int) ([]domain.Feed, error)

	UpdateFeed(ctx context.Context, f domain.Feed) error
	TouchPolled(ctx context.Context, id string, t time.Time) error

	DeleteFeed(ctx context.Context, name string) error
}

type FeedsRepository struct {
	db *sql.DB
}

func NewFeedsRepository(db *sql.DB) *FeedsRepository {
	return &FeedsRepository{db: db}
}

// Insert new feed
func (r *FeedsRepository) InsertFeed(ctx context.Context, name, url string) (domain.Feed, error) {
	const q = `
INSERT INTO feeds (name, url)
VALUES ($1, $2)
RETURNING id, created_at, updated_at, name, url,
          last_polled_at, last_changed_at, etag, last_modified;
`
	var f domain.Feed
	err := r.db.QueryRowContext(ctx, q, name, url).Scan(
		&f.ID, &f.CreatedAt, &f.UpdatedAt, &f.Name, &f.URL,
		&f.LastPolledAt, &f.LastChangeAt, &f.ETag, &f.LastModified,
	)
	return f, err
}

// Get all feeds
func (r *FeedsRepository) GetAllFeeds(ctx context.Context) ([]domain.Feed, error) {
	const q = `
SELECT id, created_at, updated_at, name, url,
       last_polled_at, last_changed_at, etag, last_modified
FROM feeds
ORDER BY created_at ASC;
`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Feed
	for rows.Next() {
		var f domain.Feed
		if err := rows.Scan(
			&f.ID, &f.CreatedAt, &f.UpdatedAt, &f.Name, &f.URL,
			&f.LastPolledAt, &f.LastChangeAt, &f.ETag, &f.LastModified,
		); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

// Get one feed by name
func (r *FeedsRepository) GetFeedByName(ctx context.Context, name string) (domain.Feed, error) {
	const q = `
SELECT id, created_at, updated_at, name, url,
       last_polled_at, last_changed_at, etag, last_modified
FROM feeds
WHERE name = $1;
`
	var f domain.Feed
	err := r.db.QueryRowContext(ctx, q, name).Scan(
		&f.ID, &f.CreatedAt, &f.UpdatedAt, &f.Name, &f.URL,
		&f.LastPolledAt, &f.LastChangeAt, &f.ETag, &f.LastModified,
	)
	return f, err
}

// Get N oldest (by last_polled_at)
func (r *FeedsRepository) GetOldestFeeds(ctx context.Context, limit int) ([]domain.Feed, error) {
	if limit <= 0 {
		limit = 1
	}
	const q = `
SELECT id, created_at, updated_at, name, url,
       last_polled_at, last_changed_at, etag, last_modified
FROM feeds
ORDER BY last_polled_at NULLS FIRST, updated_at ASC
LIMIT $1;
`
	rows, err := r.db.QueryContext(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Feed
	for rows.Next() {
		var f domain.Feed
		if err := rows.Scan(
			&f.ID, &f.CreatedAt, &f.UpdatedAt, &f.Name, &f.URL,
			&f.LastPolledAt, &f.LastChangeAt, &f.ETag, &f.LastModified,
		); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

// Update feed
func (r *FeedsRepository) UpdateFeed(ctx context.Context, f domain.Feed) error {
	const q = `
UPDATE feeds
SET name = $2,
    url = $3,
    etag = $4,
    last_modified = $5,
    last_changed_at = $6,
    updated_at = NOW()
WHERE id = $1;
`
	_, err := r.db.ExecContext(ctx, q,
		f.ID, f.Name, f.URL, f.ETag, f.LastModified, f.LastChangeAt,
	)
	return err
}

// TouchPolled
func (r *FeedsRepository) TouchPolled(ctx context.Context, id string, t time.Time) error {
	const q = `
UPDATE feeds
SET last_polled_at = $2,
    updated_at = NOW()
WHERE id = $1;
`
	_, err := r.db.ExecContext(ctx, q, id, t)
	return err
}

// Delete feed
func (r *FeedsRepository) DeleteFeed(ctx context.Context, name string) error {
	const q = `DELETE FROM feeds WHERE name = $1;`
	_, err := r.db.ExecContext(ctx, q, name)
	return err
}
