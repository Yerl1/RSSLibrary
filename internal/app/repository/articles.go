package repository

import (
	"context"
	"database/sql"
	"rsslibrary/internal/app/domain"
)

type ArticlesRepositoryInterface interface {
	GetArticles(ctx context.Context, num int, name string) ([]domain.Article, error)
}
type ArticlesRepository struct {
	db *sql.DB
}

func NewArticlesRepository(db *sql.DB) *ArticlesRepository {
	return &ArticlesRepository{db: db}
}
func (r *ArticlesRepository) GetArticles(ctx context.Context, num int, feedName string) ([]domain.Article, error) {
	const query = `
SELECT a.published_at, a.title, a.link
FROM articles a
JOIN feeds f ON f.id = a.feed_id
WHERE f.name = $1
ORDER BY a.published_at DESC NULLS LAST, a.created_at DESC
LIMIT $2;
`
	rows, err := r.db.QueryContext(ctx, query, feedName, num)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Article
	for rows.Next() {
		var a domain.Article
		if err := rows.Scan(&a.PublishedAt, &a.Title, &a.Link); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
