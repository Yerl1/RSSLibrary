package workerpool

import (
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
		return
	}
	results <- domain.Result{FeedID: job.FeedID, Articles: articles}
}
