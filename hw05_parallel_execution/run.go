package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded       = errors.New("errors limit exceeded")
	ErrZeroOrNegativeWorkerCount = errors.New("worker count cannot be zero or less")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tchan := make(chan Task, len(tasks))
	if n <= 0 {
		return ErrZeroOrNegativeWorkerCount
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(n)
	errorsCount := int32(m)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			worker(ctx, cancel, &errorsCount, tchan)
		}()
	}

	for _, t := range tasks {
		tchan <- t
	}

	close(tchan)
	wg.Wait()

	if ctx.Err() != nil {
		cancel()
		return ErrErrorsLimitExceeded
	}
	cancel()
	return nil
}

func worker(ctx context.Context, cancel context.CancelFunc, errorsCount *int32, tasks chan Task) {
	for {
		select {
		case t, ok := <-tasks:
			if !ok {
				return
			}

			if atomic.LoadInt32(errorsCount) <= 0 {
				cancel()
				return
			}

			err := t()
			if err != nil {
				atomic.AddInt32(errorsCount, -1)
				if atomic.LoadInt32(errorsCount) <= 0 {
					cancel()
					return
				}
			}

		case <-ctx.Done():
			return
		}
	}
}
