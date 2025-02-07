package storage

import (
	"time"

	domain "github.com/futig/task-scheduler/internal/domain"
)

type TaskStorage interface {
	CreateTask(task domain.Task) error
	CreateTasks(tasks []domain.Task) error
	GetTaskById(id string) (domain.Task, error)
	GetTasksByDayAndUser(day domain.DayOfTheWeek, user string) ([]domain.Task, error)
	GetCurrnetTask(user string) ([]domain.Task, error)
	UpdateTaskById(id string) error
	DeleteTaskById(id string) error
	DeleteTasksByDay(day domain.DayOfTheWeek, user string) error
} 

type TimingStorage interface {
	CreateTiming(timing domain.Timing) error
	CreateTimings(timing []domain.Timing) error
	GetTimingById(id string) (domain.Task, error)
	GetTimingsWithOffset(offset int, timeMax time.Time) ([]domain.Task, error)
	DeleteTimingById(id string) error
	DeleteTimingsByTaskId(id string) error
} 