package service

import (
	"time"

	"github.com/futig/task-scheduler/internal/domain"
	"github.com/futig/task-scheduler/internal/storage"
	"github.com/google/uuid"
)

func GetCurrentSchedule(storage storage.Storage, chatID int64) (string, bool, error) {
	return GetScheduleByWeekday(storage, chatID, time.Now().Local().Weekday())
}

func GetScheduleByWeekday(storage storage.Storage, chatID int64, weekday time.Weekday) (string, bool, error) {
	return "Schedule: Aboba", true, nil
}

func GetCurrentTask(storage storage.Storage, chatID int64) (domain.Task, bool, error) {
	return domain.Task{}, true, nil
}

func GetTaskById(storage storage.Storage, id uuid.UUID) (domain.Task, bool, error) {
	return domain.Task{}, true, nil
}

func GetTaskByPosition(storage storage.Storage, id int) (domain.Task, bool, error) {
	return domain.Task{}, true, nil
}

func CreateSchedule(storage storage.Storage, chatID int64, weekday time.Weekday) error {
	return nil
}

func UpdateSchedule(storage storage.Storage, chatID int64, weekday time.Weekday) error {
	return nil
}

func DeleteSchedule(storage storage.Storage, chatID int64, weekday time.Weekday) error {
	return nil
}

func DeleteScheduleItem(storage storage.Storage, chatID int64, weekday time.Weekday, id uuid.UUID) error {
	return nil
}
