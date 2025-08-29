package repository

import "database/sql"

type Repository struct {
	Articles ArticlesRepositoryInterface
	Feeds    FeedsRepositoryInterface
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Articles: NewArticlesRepository(db),
		Feeds:    NewFeedsRepository(db),
	}
}
