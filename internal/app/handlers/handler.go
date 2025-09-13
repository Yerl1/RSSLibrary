package handlers

import (
	"context"
	"fmt"
	"net"

	"rsslibrary/internal/app/service"
)

type Handler interface {
	Fetch(ctx context.Context, conn net.Conn)
	AddFeed(ctx context.Context, name, url string, conn net.Conn)
	SetInterval(ctx context.Context, interval string, conn net.Conn)
	SetWorkers(ctx context.Context, number string, conn net.Conn)
	ListFeeds(ctx context.Context, num int, conn net.Conn)
	DeleteFeed(ctx context.Context, name string, conn net.Conn)
	ListArticles(ctx context.Context, feedName string, num int, conn net.Conn)
}

type RequestHandler struct {
	srv *service.Service
}

func NewRequestHandler(srv *service.Service) Handler {
	return &RequestHandler{srv: srv}
}

func (h *RequestHandler) Fetch(ctx context.Context, conn net.Conn) {
	msg, err := h.srv.Fetch(ctx)
	if err != nil {
		conn.Write([]byte("Error: " + err.Error()))
		return
	}
	conn.Write([]byte(msg))
}

func (h *RequestHandler) AddFeed(ctx context.Context, name, url string, conn net.Conn) {
	feed, err := h.srv.AddFeed(ctx, name, url)
	if err != nil {
		conn.Write([]byte("Error: " + err.Error()))
		return
	}
	conn.Write([]byte(fmt.Sprintf("Feed added: %s (%s)\n", feed.Name, feed.URL)))
}

func (h *RequestHandler) SetInterval(ctx context.Context, interval string, conn net.Conn) {
	msg, err := h.srv.SetInterval(interval, ctx)
	if err != nil {
		conn.Write([]byte("Error: " + err.Error()))
		return
	}
	conn.Write([]byte(msg))
}

func (h *RequestHandler) SetWorkers(ctx context.Context, number string, conn net.Conn) {
	msg, err := h.srv.SetWorkers(number, ctx)
	if err != nil {
		conn.Write([]byte("Error: " + err.Error()))
		return
	}
	conn.Write([]byte(msg))
}
func (h *RequestHandler) ListFeeds(ctx context.Context, num int, conn net.Conn) {
	feeds, err := h.srv.ListFeeds(ctx)
	if err != nil {
		conn.Write([]byte("Error: " + err.Error()))
		return
	}

	if num > 0 && num < len(feeds) {
		feeds = feeds[len(feeds)-num:]
	}

	output := "# Available RSS Feeds\n\n"
	for i, f := range feeds {
		output += fmt.Sprintf("%d. Name: %s\n   URL: %s\n   Added: %s\n\n",
			i+1, f.Name, f.URL, f.CreatedAt.Format("2006-01-02 15:04"))
	}
	conn.Write([]byte(output))
}

func (h *RequestHandler) DeleteFeed(ctx context.Context, name string, conn net.Conn) {
	err := h.srv.DeleteFeed(ctx, name)
	if err != nil {
		conn.Write([]byte("Error: " + err.Error()))
		return
	}
	conn.Write([]byte(fmt.Sprintf("Feed '%s' deleted successfully\n", name)))
}

func (h *RequestHandler) ListArticles(ctx context.Context, feedName string, num int, conn net.Conn) {
	if num <= 0 {
		num = 3
	}
	articles, err := h.srv.GetArticles(ctx, num, feedName)
	if err != nil {
		conn.Write([]byte("Error: " + err.Error()))
		return
	}

	output := fmt.Sprintf("Feed: %s\n\n", feedName)
	for i, a := range articles {
		output += fmt.Sprintf("%d. [%s] %s\n   %s\n\n",
			i+1,
			a.PublishedAt.Format("2006-01-02"),
			a.Title,
			a.Link,
		)
	}
	conn.Write([]byte(output))
}
