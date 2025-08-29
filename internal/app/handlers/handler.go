package handlers

import (
	"context"
	"errors"

	"rsslibrary/internal/app/service"
)

type Handler interface {
	Fetch(ctx context.Context) (string, error)
	AddFeed(ctx context.Context) error
	SetInterval(ctx context.Context) (string, error)
	SetWorkers(ctx context.Context) (string, error)
	List(ctx context.Context) (string, error)
	DeleteFeed(ctx context.Context) error
	Articles(ctx context.Context) (string, error)
}

type RequestHandler struct {
	srv         *service.Service
	fetchStatus bool
}

func NewRequestHandler(srv *service.Service) *RequestHandler {
	return &RequestHandler{fetchStatus: false, srv: srv}
}

func (this *RequestHandler) Fetch(ctx context.Context) (string, error) {
	if this.fetchStatus {
		return "", errors.New("Background process is already running")
	}
	this.fetchStatus = true

	return "Background process for fetching feed has started (interval = 3 minutes, workers = 3)", nil
}

func (this *RequestHandler) AddFeed(ctx context.Context) error {
	return nil
}

func (this *RequestHandler) SetInterval(ctx context.Context) (string, error) {
	if !this.fetchStatus {
		return "", errors.New("Fetch process is not started yet")
	}
	return "Interval of fetching feeds changed from 3 minutes to 2 minutes", nil
}

func (this *RequestHandler) SetWorkers(ctx context.Context) (string, error) {
	if !this.fetchStatus {
		return "", errors.New("Fetch process is not started yet")
	}
	return "Number of workers changed from 3 to 5", nil
}

func (this *RequestHandler) List(ctx context.Context) (string, error) {
	return "", nil
}

func (this *RequestHandler) DeleteFeed(ctx context.Context) error {
	return nil
}

func (this *RequestHandler) Articles(ctx context.Context) (string, error) {
	return "", nil
}
