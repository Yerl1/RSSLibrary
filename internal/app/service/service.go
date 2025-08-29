package service

import (
	"sync"

	"rsslibrary/internal/app/service/workerpool"
)

type Service interface {
	Fetch()
}

type RequestServer struct{}

func (this *RequestServer) Fetch() {
	bufferSize := 1000
	var wg sync.WaitGroup
	dispatcher := workerpool.NewDispatcher(bufferSize, &wg, 500)
	_ = dispatcher
}
