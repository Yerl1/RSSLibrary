package domain

import "time"

type Article struct {
	ID          string    `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	Title       string    `db:"title"`
	Link        string    `db:"link"`
	PublishedAt time.Time `db:"published_at"`
	Description string    `db:"description"`
	FeedID      string    `db:"feed_id"`
}
