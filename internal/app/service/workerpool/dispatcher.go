package workerpool

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type Dispatcher interface {
	AddWorker(w WorkerLauncher)
	RemoveWorker(minWorkers int)
	LaunchWorker(id int, w WorkerLauncher)
	ScaleWorkers(minWorkers, maxWorkers, loadThreshold int)
	MakeRequest(r string)
	Stop(ctx context.Context)
	StartDispatcher(ctx context.Context)
}

type dispatcher struct {
	inCh            chan string
	wg              *sync.WaitGroup
	mu              sync.Mutex
	ticker          *time.Ticker
	workerCount     int
	minWorkerNumber int
	stopCh          chan struct{}
}

func NewDispatcher(b int, wg *sync.WaitGroup, ticker *time.Ticker) Dispatcher {
	minWorkerNumber, _ := strconv.Atoi(os.Getenv("CLI_APP_WORKERS_COUNT"))
	return &dispatcher{
		inCh:            make(chan string, b),
		wg:              wg,
		stopCh:          make(chan struct{}, 50),
		ticker:          ticker,
		minWorkerNumber: minWorkerNumber,
	}
}

func (d *dispatcher) StartDispatcher(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				d.Stop(ctx)
				return
			case <-d.ticker.C:
				for i := 0; i < d.minWorkerNumber; i++ {
					fmt.Printf("Starting worker with the id %d\n", i)
					w := &Worker{
						Id: i,
						Wg: d.wg,
					}
					d.AddWorker(w)
				}
			}
		}
	}()
}

func (d *dispatcher) AddWorker(w WorkerLauncher) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.workerCount++
	d.wg.Add(1)
	w.LaunchWorker(d.inCh, d.stopCh)
}

func (d *dispatcher) RemoveWorker(minWorkers int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.workerCount > minWorkers {
		d.workerCount--
		d.stopCh <- struct{}{}
	}
}

func (d *dispatcher) ScaleWorkers(minWorkers, maxWorkers, loadThreshold int) {
	ticker := time.NewTicker(time.Microsecond)
	defer ticker.Stop()

	for range ticker.C {
		load := len(d.inCh)
		if load > loadThreshold && d.workerCount < maxWorkers {
			fmt.Println("Scaling Triggered")
			newWorker := &Worker{
				Wg: d.wg,
				Id: d.workerCount,
			}
			d.AddWorker(newWorker)
		} else if float64(load) < float64(0.75)*float64(loadThreshold) && d.workerCount > minWorkers {
			fmt.Println("Reducing Triggered")
			d.RemoveWorker(minWorkers)
		}
	}
}

func (d *dispatcher) LaunchWorker(id int, w WorkerLauncher) {
	w.LaunchWorker(d.inCh, d.stopCh)
	d.mu.Lock()
	d.workerCount++
	d.mu.Unlock()
}

func (d *dispatcher) MakeRequest(r string) {
	select {
	case d.inCh <- r:
	default:
		fmt.Println("Request channel is full. Dropping request.")
	}
}

func (d *dispatcher) Stop(ctx context.Context) {
	fmt.Println("\nstop called")
	close(d.inCh)
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
		for i := 0; i < d.workerCount; i++ {
			d.stopCh <- struct{}{}
		}
	}
	d.wg.Wait()
}
