package tempstorage

import (
	"fmt"
	"sync"
	"time"

	"github.com/futig/task-scheduler/internal/domain"
	t "github.com/futig/task-scheduler/pkg/time"

	"github.com/google/uuid"
)

type TempTaskStorage struct {
	Storage sync.Map
}

func (s *TempStorageContext) CreateTask(task domain.Task) error {
	if _, ok := s.Tasks.Load(task.Id); ok {
		return fmt.Errorf("задача с таким id уже существует %s", task.Id.String())
	}
	s.Tasks.Store(task.Id, task)
	return nil
}

func (s *TempStorageContext) CreateTasks(tasks []domain.Task) error {
	for _, task := range tasks {
		if _, ok := s.Tasks.Load(task.Id); ok {
			return fmt.Errorf("задача с таким id уже существует %s", task.Id.String())
		}
		s.Tasks.Store(task.Id, task)
	}
	return nil
}

func (s *TempStorageContext) GetTaskById(id uuid.UUID) (domain.Task, bool, error) {
	if val, ok := s.Tasks.Load(id); ok {
		return val.(domain.Task), ok, nil
	}
	return domain.Task{}, false, fmt.Errorf("задачи с таким id не существует: %s", id.String())
}

func (s *TempStorageContext) GetTasksByDayAndUser(day time.Weekday, chatID int64) ([]domain.Task, error) {
	result := make([]domain.Task, 100)
	for _, value := range s.Tasks.Range {
		task := value.(domain.Task)
		if task.Day == day && task.ChatId == chatID {
			result = append(result, task)
		}
	}
	return result, nil
}

func (s *TempStorageContext) GetCurrnetTask(chatID int64) (domain.Task, bool, error) {
	for _, value := range s.Tasks.Range {
		task := value.(domain.Task)
		day := time.Now().Weekday()
		curTime := t.CurrentTimeToMinutes()
		if task.Day == day && task.ChatId == chatID && task.Start <= curTime && task.End > curTime {
			return task, true, nil
		}
	}
	return domain.Task{}, false, nil
}

func (s *TempStorageContext) UpdateTaskById(id uuid.UUID, task domain.Task) error {
	if _, ok := s.Tasks.Load(id); ok {
		s.Tasks.Store(id, task)
		return nil
	}
	return fmt.Errorf("задачи с таким id не существует: %s", id.String())
}

func (s *TempStorageContext) DeleteTaskById(id uuid.UUID) error {
	s.Tasks.Delete(id)
	return nil
}

func (s *TempStorageContext) DeleteTasksByDay(day time.Weekday, chatID int64) error {
	for _, value := range s.Tasks.Range {
		task := value.(domain.Task)
		day := time.Now().Weekday()
		if task.Day == day && task.ChatId == chatID {
			s.Tasks.Delete(task.Id)
		}
	}
	return nil
}
