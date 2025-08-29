package repository

import "database/sql"

type FeedsRepositoryInterface interface {
	// Define methods for FeedsRepositoryInterface
}

type FeedsRepository struct {
	db *sql.DB
}

func NewFeedsRepository(db *sql.DB) *FeedsRepository {
	return &FeedsRepository{db: db}
}
