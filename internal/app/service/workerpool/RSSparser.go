package workerpool

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"

	"rsslibrary/internal/app/domain"
)

type rssFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items       []rssItem `xml:"item"`
	} `xml:"channel"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchRSS(feedID string, url string) ([]domain.Article, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rss rssFeed
	if err := xml.Unmarshal(data, &rss); err != nil {
		return nil, err
	}

	var articles []domain.Article
	now := time.Now()

	for _, item := range rss.Channel.Items {
		var published time.Time
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			published = t
		} else if t, err := time.Parse(time.RFC1123, item.PubDate); err == nil {
			published = t
		} else {
			published = now
		}

		articles = append(articles, domain.Article{
			CreatedAt:   now,
			UpdatedAt:   now,
			Title:       item.Title,
			Link:        item.Link,
			PublishedAt: published,
			Description: item.Description,
			FeedID:      feedID,
		})
	}

	return articles, nil
}
