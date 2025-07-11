package periodic_tasks

import (
	"context"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Task interface {
	RunTask() error
}

// RunMonthlyTask - Позволяет запустить периодическую задачу, выполняемую каждое 2 число месяца в 10:00 UTC+0
func RunMonthlyTask[T Task](ctx context.Context, wg *sync.WaitGroup, log *zap.Logger, task T) error {
	defer wg.Done()
	for {
		now := time.Now()
		next := now.AddDate(0, 1, 0)

		nextDate := time.Date(next.Year(), next.Month(), 2, 10, 0, 0, 0, time.UTC)

		select {
		case <-time.After(nextDate.Sub(now)):
			if err := task.RunTask(); err != nil {
				log.Error("RunMonthlyTask", zap.Error(err))
			}
		case <-ctx.Done():
			log.Info("RunMonthlyTask cancelled")
			return ctx.Err()
		}
	}
}
