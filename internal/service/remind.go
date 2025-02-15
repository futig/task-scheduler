package service

import (
	"time"

	"github.com/futig/task-scheduler/internal/domain"
	"github.com/futig/task-scheduler/internal/storage"
)

func GetRemindsForPeriod(storage storage.Storage, period time.Duration) ([]domain.TaskRemind, error) {
	totalMinutes := int(period.Minutes())
	reminds, err := storage.GetRemindsForPeriod(totalMinutes)
	return reminds, err
}
