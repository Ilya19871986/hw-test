package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrorsIllegalArgument = errors.New("errors illegal arguments")

type Task func() error

func Run(tasks []Task, n int, m int) error {
	if len(tasks) == 0 || m <= 0 || n <= 0 {
		return ErrorsIllegalArgument
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	taskCh := make(chan Task, len(tasks))
	var errorCount int32
	var wg sync.WaitGroup
	// Запускаем воркеров.
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskCh {
				if err := task(); err != nil {
					if atomic.AddInt32(&errorCount, 1) >= int32(m) {
						cancel()
					}
				}
				select {
				case <-ctx.Done():
					return
				default:
				}
			}
		}()
	}
	// Отправки задач в канал для обработки воркерами.
sendLoop:
	for _, task := range tasks {
		if atomic.LoadInt32(&errorCount) >= int32(m) {
			break sendLoop
		}

		select {
		case taskCh <- task:
		case <-ctx.Done():
			break sendLoop
		}
	}
	close(taskCh)

	wg.Wait()

	if atomic.LoadInt32(&errorCount) >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
