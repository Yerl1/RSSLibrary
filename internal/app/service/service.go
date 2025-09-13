package service

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"rsslibrary/internal/app/domain"
	"rsslibrary/internal/app/repository"
	"rsslibrary/internal/app/service/workerpool"
)

type ServiceInterface interface {
	Fetch(ctx context.Context) (string, error)
	SetInterval(interval string, ctx context.Context) (string, error)
	SetWorkers(number string, ctx context.Context) (string, error)
	AddFeed(ctx context.Context, name, url string) (domain.Feed, error)
	ListFeeds(ctx context.Context) ([]domain.Feed, error)
	DeleteFeed(ctx context.Context, name string) error
	GetArticles(ctx context.Context, num int, feedName string) ([]domain.Article, error)
}

type Service struct {
	repository   *repository.Repository
	FETCH_STATUS bool
	dispatcher   workerpool.Dispatcher
}

func NewService(repo *repository.Repository) *Service {
	var wg sync.WaitGroup
	workers, _ := strconv.Atoi(os.Getenv("CLI_APP_WORKERS_COUNT"))
	return &Service{
		repository: repo,
		dispatcher: workerpool.NewDispatcher(50000, &wg, workers),
	}
}

func (s *Service) Fetch(ctx context.Context) (string, error) {
	if s.FETCH_STATUS {
		return "Background process is already running", nil
	}

	minutes, _ := strconv.Atoi(os.Getenv("CLI_APP_TIMER_INTERVAL")[:len(os.Getenv("CLI_APP_TIMER_INTERVAL"))-1])
	s.dispatcher.SetInterval(time.Minute * time.Duration(minutes))

	jobGenerator := func() []domain.Job {
		feeds, err := s.repository.Feeds.GetAllFeeds(ctx)
		if err != nil {
			return nil
		}
		jobs := make([]domain.Job, 0, len(feeds))
		for _, f := range feeds {
			jobs = append(jobs, domain.Job{FeedID: f.ID, URL: f.URL})
		}
		return jobs
	}

	s.dispatcher.StartDispatcher(ctx, jobGenerator)

	go func() {
		for res := range s.dispatcher.Results() {
			if res.Err != nil {
				fmt.Printf("Error processing feed %s: %v\n", res.FeedID, res.Err)
				continue
			}
			if len(res.Articles) > 0 {
				if err := s.repository.Articles.InsertArticles(ctx, res.Articles); err != nil {
					fmt.Printf("DB insert error for feed %s: %v\n", res.FeedID, err)
				}
			}
			_ = s.repository.Feeds.TouchPolled(ctx, res.FeedID, time.Now())
		}
	}()

	s.FETCH_STATUS = true
	return fmt.Sprintf("Background process started (interval=%s, workers=%s)",
		os.Getenv("CLI_APP_TIMER_INTERVAL"),
		os.Getenv("CLI_APP_WORKERS_COUNT")), nil
}

func (s *Service) SetInterval(interval string, ctx context.Context) (string, error) {
	if !s.FETCH_STATUS {
		return "Background process is not running", nil
	}
	size := len(interval)
	minutes, err := strconv.Atoi(interval[:size-1])
	if err != nil {
		return "", err
	}

	previous := s.dispatcher.GetInterval()
	current := time.Duration(minutes) * time.Minute
	s.dispatcher.SetInterval(current)

	message := fmt.Sprintf(
		"Interval of fetching feeds changed from %d minutes to %d minutes",
		int(previous.Minutes()),
		minutes,
	)
	return message, nil
}

func (s *Service) SetWorkers(number string, ctx context.Context) (string, error) {
	if !s.FETCH_STATUS {
		return "Background process is not running", nil
	}

	workersNum, err := strconv.Atoi(number)
	if err != nil {
		return "", err
	}

	previous := s.dispatcher.GetWorkerCount()
	s.dispatcher.SetWorkers(workersNum)

	message := fmt.Sprintf(
		"Number of workers changed from %d to %d",
		previous,
		workersNum,
	)
	return message, nil
}

func (s *Service) AddFeed(ctx context.Context, name, url string) (domain.Feed, error) {
	return s.repository.Feeds.InsertFeed(ctx, name, url)
}

func (s *Service) ListFeeds(ctx context.Context) ([]domain.Feed, error) {
	return s.repository.Feeds.GetAllFeeds(ctx)
}

func (s *Service) DeleteFeed(ctx context.Context, name string) error {
	return s.repository.Feeds.DeleteFeed(ctx, name)
}

func (s *Service) GetArticles(ctx context.Context, num int, feedName string) ([]domain.Article, error) {
	return s.repository.Articles.GetArticles(ctx, num, feedName)
}
