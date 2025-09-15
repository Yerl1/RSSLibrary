package workerpool

import (
	"fmt"
	"log"
	"rsslibrary/internal/app/domain"
	"sync"
)

type Worker struct {
	Id int
	Wg *sync.WaitGroup
}

func (w *Worker) LaunchWorker(jobs <-chan domain.Job, results chan<- domain.Result, stopCh <-chan struct{}) {
	go func() {
		defer w.Wg.Done()
		for {
			select {
			case job, ok := <-jobs:
				if !ok {
					return
				}
				fmt.Println("I have started worker")
				fmt.Println("Working on ", job)
				w.ProcessJob(job, results)
			case <-stopCh:
				return
			}
		}
	}()
}

func (w *Worker) ProcessJob(job domain.Job, results chan<- domain.Result) {
	articles, err := FetchRSS(job.FeedID, job.URL)
	if err != nil {
		results <- domain.Result{FeedID: job.FeedID, Err: err}
		log.Println("got this error in ProcessJob", err.Error())
		return
	}
	results <- domain.Result{FeedID: job.FeedID, Articles: articles}
}
