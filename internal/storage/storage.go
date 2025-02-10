package storage

import (
	"time"

	"github.com/futig/task-scheduler/internal/domain"

	"github.com/google/uuid"
)

type TaskStorage interface {
	CreateTask(task domain.Task) error
	CreateTasks(tasks []domain.Task) error
	GetTaskById(id uuid.UUID) (domain.Task, bool, error)
	GetTasksByDayAndUser(day time.Weekday, chatID int64) ([]domain.Task, error)
	GetCurrnetTask(chatID int64) (domain.Task, bool, error)
	UpdateTaskById(id uuid.UUID, task domain.Task) error
	DeleteTaskById(id uuid.UUID) error
	DeleteTasksByDay(day time.Weekday, chatID int64) error
}

type RemindsStorage interface {
	CreateRemind(remind domain.Remind) error
	CreateReminds(remind []domain.Remind) error
	GetRemindById(id uuid.UUID) (domain.TaskRemind, bool, error)
	GetRemindsForPeriod(period int) ([]domain.TaskRemind, error) // period в минутах на будущее
	DeleteRemindById(id uuid.UUID) error
	DeleteRemindsByTaskId(id uuid.UUID) error
}
