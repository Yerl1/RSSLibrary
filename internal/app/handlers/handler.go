package handlers

import "context"

type Handler interface {
	Fetch(ctx context.Context)
	AddFeed(ctx context.Context)
	SetInterval(ctx context.Context)
	SetWorkers(ctx context.Context)
	List(ctx context.Context)
	DeleteFeed(ctx context.Context)
	Articles(ctx context.Context)
}
