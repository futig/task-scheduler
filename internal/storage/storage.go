package storage

import (
	"time"

	"github.com/futig/task-scheduler/internal/domain"
	"github.com/futig/task-scheduler/internal/domain/enums"
)

type TaskStorage interface {
	CreateTask(task domain.Task) error
	CreateTasks(tasks []domain.Task) error
	GetTaskById(id string) (domain.Task, error)
	GetTasksByDayAndUser(day enums.DayOfTheWeek, user string) ([]domain.Task, error)
	GetCurrnetTask(user string) ([]domain.Task, error)
	UpdateTaskById(id string) error
	DeleteTaskById(id string) error
	DeleteTasksByDay(day enums.DayOfTheWeek, user string) error
}

type RemindsStorage interface {
	CreateTiming(timing domain.Remind) error
	CreateTimings(timing []domain.Remind) error
	GetTimingById(id string) (domain.Task, error)
	GetTimingsWithOffset(offset int, timeMax time.Time) ([]domain.Task, error)
	DeleteTimingById(id string) error
	DeleteTimingsByTaskId(id string) error
}
