package handlers

import (
	"context"
	"fmt"
	"net"

	"rsslibrary/internal/app/service"
)

type Handler interface {
	Fetch(ctx context.Context)
	AddFeed(ctx context.Context)
	SetInterval(ctx context.Context)
	SetWorkers(ctx context.Context)
	List(ctx context.Context)
	DeleteFeed(ctx context.Context)
	Articles(ctx context.Context)
}

type RequestHandler struct {
	srv *service.Service
}

func NewRequestHandler(srv *service.Service) *RequestHandler {
	return &RequestHandler{srv: srv}
}

func (this *RequestHandler) Fetch(ctx context.Context, conn net.Conn) {
	msg, err := this.srv.Fetch(ctx)
	if err != nil {
		// write to log
		fmt.Println(err.Error())
		return
	}
	fmt.Println(msg)
	conn.Write([]byte(msg))
}

func (this *RequestHandler) AddFeed(ctx context.Context) error {
	return nil
}

func (this *RequestHandler) SetInterval(ctx context.Context) (string, error) {
	return "Interval of fetching feeds changed from 3 minutes to 2 minutes", nil
}

func (this *RequestHandler) SetWorkers(ctx context.Context) (string, error) {
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
