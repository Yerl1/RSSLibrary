package repository

import (
	"context"
	"database/sql"
)

type FeedsRepositoryInterface interface {
	// Define methods for FeedsRepositoryInterface
}

type FeedsRepository struct {
	db *sql.DB
}

func NewFeedsRepository(db *sql.DB) *FeedsRepository {
	return &FeedsRepository{db: db}
}

func (r *FeedsRepository) InsertFeed(ctx context.Context) error {
	// Implement the method to insert a feed into the database
	return nil
}

func (r *FeedsRepository) GetAllFeeds(ctx context.Context) ([]string, error) {
	query := "SELECT url FROM feeds ORDER BY updated_at ASC"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var feeds []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		feeds = append(feeds, url)
	}
	return feeds, nil
}
