package domain

type Feed struct {
	Id             string
	CreatedAt      string
	UpdatedAt      string
	Name           string
	Url            string
	Last_polled_at string
	Last_change_at string
	Etag           string
	Last_modified  string
}
