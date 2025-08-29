package workerpool

import (
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	Id int
	Wg *sync.WaitGroup
}

func (w *Worker) LaunchWorker(in chan string, stopCh chan struct{}) {
	go func() {
		defer w.Wg.Done()
		for {
			select {
			case req, open := <-in:
				if !open {
					fmt.Println("Stop worker:", w.Id, " Reason: request channel is closed")
					return
				}
				w.ProcessRequest(req)
				time.Sleep(1 * time.Microsecond)
			case <-stopCh:
				fmt.Println("Stopping worker:", w.Id, " Reason: worker was intentionally removed")
				return
			}
		}
	}()
}

func (w *Worker) ProcessRequest(req string) {
}
