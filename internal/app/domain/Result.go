package domain

type Result struct {
	FeedID   string
	Articles []Article
	Err      error
}
