package worker_pool_pattern

import (
	"context"
	"sync"
)

type Job interface {
	ExecFn(ctx context.Context) (interface{}, error)
}

type WorkerPoolOption func(*WorkerPool)

type WorkerPool struct {
	workerCounter int
	jobChan       chan Job
	resultChan    chan Result
	isFetchResult bool
}

type Result struct {
	Value interface{}
	Err   error
}

func WithIsFetchResult(ok bool) WorkerPoolOption {
	return func(wp *WorkerPool) {
		wp.isFetchResult = ok
	}
}

func New(workerCounter int, chanBuffer int, opts ...WorkerPoolOption) *WorkerPool {
	wp := &WorkerPool{
		workerCounter: workerCounter,
		jobChan:       make(chan Job, chanBuffer),
		resultChan:    make(chan Result, chanBuffer),
		isFetchResult: false,
	}
	for _, opt := range opts {
		opt(wp)
	}
	return wp
}

func (wp *WorkerPool) Run(ctx context.Context) error {
	defer close(wp.resultChan)
	var wg sync.WaitGroup

	for idx := 0; idx < wp.workerCounter; idx++ {
		wg.Add(1)
		worker := func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-wp.jobChan:
					if !ok {
						return
					}
					result, err := job.ExecFn(ctx)
					if !wp.isFetchResult {
						continue
					}
					wp.resultChan <- Result{
						Value: result,
						Err:   err,
					}
				}
			}
		}
		go worker()
	}
	wg.Wait()
	return nil
}

func (wp *WorkerPool) Results() <-chan Result {
	return wp.resultChan
}

func (wp *WorkerPool) Push(job Job) {
	wp.jobChan <- job
}
