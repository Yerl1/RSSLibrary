package repository

import (
	"database/sql"
)

type ArticlesRepositoryInterface interface {
	// Define methods for ArticlesRepositoryInterface
}
type ArticlesRepository struct {
	db *sql.DB
}

func NewArticlesRepository(db *sql.DB) *ArticlesRepository {
	return &ArticlesRepository{db: db}
}

// func (r *ArticlesRepository) InsertArticle(ctx context.Context, ) error {

// }
