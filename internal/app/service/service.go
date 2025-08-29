package service

import (
	"context"
	"errors"
	"os"
	"strconv"
	"sync"
	"time"

	"rsslibrary/internal/app/repository"
	"rsslibrary/internal/app/service/workerpool"
)

type ServiceInterface interface {
	Fetch()
}

type Service struct {
	repository   *repository.Repository
	FETCH_STATUS bool
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		repository: repo,
	}
}

func (this *Service) Fetch(ctx context.Context) (string, error) {
	if this.FETCH_STATUS {
		return "Background process is already running", errors.New("Noooo")
	}
	var wg sync.WaitGroup
	size := len(os.Getenv("CLI_APP_TIMER_INTERVAL"))
	minutes, err := strconv.Atoi(os.Getenv("CLI_APP_TIMER_INTERVAL")[:size-1])
	if err != nil {
		return "", err
	}
	ticker := time.NewTicker(time.Minute * time.Duration(minutes))
	dispatcher := workerpool.NewDispatcher(50000, &wg, ticker)
	dispatcher.StartDispatcher(ctx)
	this.FETCH_STATUS = true
	message := "The background process for fetching feeds has started "
	message += "(interval  = " + os.Getenv("CLI_APP_TIMER_INTERVAL") + ", workers = " + os.Getenv("CLI_APP_WORKERS_COUNT") + ")"
	return message, nil
}
