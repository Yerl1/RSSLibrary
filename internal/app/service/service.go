package service

import (
	"sync"

	"rsslibrary/internal/app/repository"
	"rsslibrary/internal/app/service/workerpool"
)

type ServiceInterface interface {
	Fetch()
}

type Service struct {
	repository *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		repository: repo,
	}
}

func (this *Service) Fetch() {
	bufferSize := 1000
	var wg sync.WaitGroup
	dispatcher := workerpool.NewDispatcher(bufferSize, &wg, 500)
	_ = dispatcher
}
