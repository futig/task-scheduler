package app

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/futig/task-scheduler/internal/service"
)

func reminderManager(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	timerCh := time.After(cfg.RemindsCheckInterval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-timerCh:
			tasks, err := service.GetRemindsForPeriod(wCfg.Storage, cfg.RemindsCheckInterval)
			if err != nil {
				log.Print(err.Error())
				continue
			}
			for _, task := range tasks {
				wCfg.RemindsCh <- task
			}
		}
	}
}
