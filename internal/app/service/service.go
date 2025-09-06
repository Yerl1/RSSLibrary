package service

import (
	"context"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"rsslibrary/internal/app/repository"
	"rsslibrary/internal/app/service/workerpool"
)

type ServiceInterface interface {
	Fetch()
	SetInterval()
}

type Service struct {
	repository   *repository.Repository
	FETCH_STATUS bool
	dispatcher   workerpool.Dispatcher
}

func NewService(repo *repository.Repository) *Service {
	var wg sync.WaitGroup
	return &Service{
		repository: repo,
		dispatcher: workerpool.NewDispatcher(50000, &wg),
	}
}

func (this *Service) Fetch(ctx context.Context) (string, error) {
	if this.FETCH_STATUS {
		return "Background process is already running", nil
	}
	size := len(os.Getenv("CLI_APP_TIMER_INTERVAL"))
	minutes, err := strconv.Atoi(os.Getenv("CLI_APP_TIMER_INTERVAL")[:size-1])
	if err != nil {
		return "", err
	}
	this.dispatcher.SetInterval(time.Minute * time.Duration(minutes))
	this.dispatcher.StartDispatcher(ctx)
	this.FETCH_STATUS = true
	message := "The background process for fetching feeds has started "
	message += "(interval  = " + os.Getenv("CLI_APP_TIMER_INTERVAL") + ", workers = " + os.Getenv("CLI_APP_WORKERS_COUNT") + ")"
	return message, nil
}

func (this *Service) SetInterval(interval string, ctx context.Context) (string, error) {
	if !this.FETCH_STATUS {
		return "Background process is not running", nil
	}
	// BONUS: add checking for integer
	size := len(interval)
	minutes, err := strconv.Atoi(interval[:size-1])
	if err != nil {
		return "", err
	}
	previosInterval := this.dispatcher.GetInterval()
	currentInterval := time.Duration(minutes) * time.Minute
	this.dispatcher.SetInterval(currentInterval)

	// Function for tools.go
	getMinutesHelpfunction := func(time string) string {
		var builder strings.Builder
		for i := 0; i < len(time); i++ {
			if time[i] < '0' || time[i] > '9' {
				break
			}
			builder.WriteByte(time[i])
		}
		return builder.String()
	}

	message := "Interval of fetching feeds changed from " + getMinutesHelpfunction(previosInterval.String()) + " minutes to " + getMinutesHelpfunction(currentInterval.String()) + " minutes"
	return message, nil
}
