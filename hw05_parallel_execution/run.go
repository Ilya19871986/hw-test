package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrorsIllegalArgument  = errors.New("errors illegal arguments")
)

type Task func() error

func Run(tasks []Task, workersNum int, maxErrors int) error {
	if len(tasks) == 0 || maxErrors <= 0 || workersNum <= 0 {
		return ErrorsIllegalArgument
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	taskCh := make(chan Task, len(tasks))
	var errorCount int32
	var wg sync.WaitGroup
	// Запускаем воркеров.
	for i := 0; i < workersNum; i++ {
		wg.Add(1)
		go worker(ctx, &wg, taskCh, &errorCount, maxErrors, cancel)
	}
	// Отправки задач в канал для обработки воркерами.
sendLoop:
	for _, task := range tasks {
		if atomic.LoadInt32(&errorCount) >= int32(maxErrors) {
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
	if atomic.LoadInt32(&errorCount) >= int32(maxErrors) {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func worker(ctx context.Context, wg *sync.WaitGroup, taskCh <-chan Task,
	errorCount *int32, maxErrors int, cancel context.CancelFunc,
) {
	defer wg.Done()

	for task := range taskCh {
		if err := task(); err != nil {
			if atomic.AddInt32(errorCount, 1) >= int32(maxErrors) {
				cancel()
			}
		}
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}
