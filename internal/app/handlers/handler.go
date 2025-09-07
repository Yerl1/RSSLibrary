package handlers

import (
	"context"
	"fmt"
	"net"
	"strings"

	"rsslibrary/internal/app/service"
)

type Handler interface {
	Fetch(ctx context.Context, conn net.Conn)
	AddFeed(data string, ctx context.Context, conn net.Conn)
	SetInterval(inverval string, ctx context.Context, conn net.Conn)
	SetWorkers(number string, ctx context.Context, conn net.Conn)
	List(ctx context.Context)
	DeleteFeed(ctx context.Context)
	Articles(ctx context.Context)
}

type RequestHandler struct {
	srv *service.Service
}

func NewRequestHandler(srv *service.Service) Handler {
	return &RequestHandler{srv: srv}
}

func (this *RequestHandler) Fetch(ctx context.Context, conn net.Conn) {
	msg, err := this.srv.Fetch(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	conn.Write([]byte(msg))
}

func (this *RequestHandler) AddFeed(data string, ctx context.Context, conn net.Conn) {
	// Correct Flag parse

	// Temp quotation parse
	dataArgs := strings.Split(data, " ")
	name := dataArgs[1]
	url := dataArgs[3]
	message, err := this.srv.AddFeed(url, name, ctx)
	if err != nil {
		conn.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	conn.Write([]byte(message))
}

func (this *RequestHandler) SetInterval(inverval string, ctx context.Context, conn net.Conn) {
	msg, err := this.srv.SetInterval(inverval, ctx)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	conn.Write([]byte(msg))
}

func (this *RequestHandler) SetWorkers(number string, ctx context.Context, conn net.Conn) {
	msg, err := this.srv.SetWorkers(number, ctx)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	conn.Write([]byte(msg))
}

func (this *RequestHandler) List(ctx context.Context) {
}

func (this *RequestHandler) DeleteFeed(ctx context.Context) {
}

func (this *RequestHandler) Articles(ctx context.Context) {
}
