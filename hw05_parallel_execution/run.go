package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type Worker struct {
	errorCount *int
	tasks      chan Task
	wg         *sync.WaitGroup
	mu         *sync.RWMutex
}

func (w *Worker) Run(m int) {
	defer w.wg.Done()
	for {
		task, ok := <-w.tasks
		if !ok {
			return
		}
		w.mu.RLock()
		errCnt := *w.errorCount
		w.mu.RUnlock()
		if errCnt >= m {
			return
		}

		if err := task(); err != nil {
			w.mu.Lock()
			*w.errorCount++
			w.mu.Unlock()
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var (
		errorCount int
		wg         sync.WaitGroup
		mu         sync.RWMutex
	)
	wg.Add(n)
	taskChan := make(chan Task, len(tasks))
	for i := 0; i < n; i++ {
		worker := Worker{
			tasks:      taskChan,
			errorCount: &errorCount,
			wg:         &wg,
			mu:         &mu,
		}
		go worker.Run(m)
	}
	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)
	wg.Wait()
	if errorCount >= n {
		return ErrErrorsLimitExceeded
	}
	return nil
}
