package domain

import "time"

type Feed struct {
	ID           string     `db:"id"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	Name         string     `db:"name"`
	URL          string     `db:"url"`
	LastPolledAt *time.Time `db:"last_polled_at"`
	LastChangeAt *time.Time `db:"last_change_at"`
	ETag         *string    `db:"etag"`
	LastModified *string    `db:"last_modified"`
}
