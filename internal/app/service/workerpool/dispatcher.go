package workerpool

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"rsslibrary/internal/app/domain"
)

type Dispatcher interface {
	StartDispatcher(ctx context.Context, jobGenerator func() []domain.Job)
	SetInterval(interval time.Duration)
	GetInterval() time.Duration
	SetWorkers(number int)
	GetWorkerCount() int
	Stop(ctx context.Context)
	Results() <-chan domain.Result
}

type dispatcher struct {
	jobs             chan domain.Job
	results          chan domain.Result
	wg               *sync.WaitGroup
	mu               sync.Mutex
	ticker           *time.Ticker
	workers          []*Worker
	interval         time.Duration
	stopCh           chan struct{}
	updateIntervalCh chan time.Duration
}

func NewDispatcher(buf int, wg *sync.WaitGroup, workerCount int) Dispatcher {
	minutes, _ := strconv.Atoi(os.Getenv("CLI_APP_TIMER_INTERVAL")[:len(os.Getenv("CLI_APP_TIMER_INTERVAL"))-1])
	d := &dispatcher{
		jobs:             make(chan domain.Job, buf),
		wg:               wg,
		stopCh:           make(chan struct{}),
		updateIntervalCh: make(chan time.Duration, 1),
		ticker:           time.NewTicker(time.Duration(minutes) * time.Minute),
	}
	d.SetWorkers(workerCount)
	return d
}

func (d *dispatcher) SetInterval(interval time.Duration) {
	d.updateIntervalCh <- interval
}

func (d *dispatcher) GetInterval() time.Duration {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.interval
}

func (d *dispatcher) SetWorkers(count int) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if count < len(d.workers) {
		diff := len(d.workers) - count
		for i := 0; i < diff; i++ {
			d.stopCh <- struct{}{}
		}
		d.workers = d.workers[:count]
		return
	}

	for i := len(d.workers); i < count; i++ {
		w := &Worker{Id: i + 1, Wg: d.wg}
		d.wg.Add(1)
		w.LaunchWorker(d.jobs, d.results, d.stopCh)
		d.workers = append(d.workers, w)
	}
}

func (d *dispatcher) GetWorkerCount() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.workers)
}

func (d *dispatcher) StartDispatcher(ctx context.Context, jobGenerator func() []domain.Job) {
	log.Printf("Dispatcher started")
	go func() {
		for {
			select {
			case <-ctx.Done():
				d.Stop(ctx)
				return
			case newInterval := <-d.updateIntervalCh:
				d.ticker.Stop()
				d.ticker = time.NewTicker(newInterval)
				d.interval = newInterval
				log.Println("Interval updated")
			case <-d.ticker.C:
				log.Printf("Workers are starting to work....")
				for _, job := range jobGenerator() {
					d.jobs <- job
				}
			}
		}
	}()
}

func (d *dispatcher) Stop(ctx context.Context) {
	fmt.Println("Graceful shutdown started...")
	if d.ticker != nil {
		d.ticker.Stop()
	}
	close(d.stopCh)
	close(d.jobs)

	done := make(chan struct{})
	go func() {
		d.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("All workers stopped gracefully")
	case <-ctx.Done():
		fmt.Println("Timeout reached, forcing shutdown")
	}
}

func (d *dispatcher) Results() <-chan domain.Result {
	return d.results
}
