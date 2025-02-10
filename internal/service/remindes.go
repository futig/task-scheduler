package service

import (
	"time"

	"github.com/futig/task-scheduler/internal/domain"
	"github.com/futig/task-scheduler/internal/storage"
)

func GetRemindsForPastPeriod(storage storage.RemindsStorage, period time.Duration) ([]domain.TaskRemind, error) {
	return nil, nil
}
