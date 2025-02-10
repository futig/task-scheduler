package tempstorage

import (
	"fmt"
	"sync"
	"time"

	"github.com/futig/task-scheduler/internal/domain"
	"github.com/google/uuid"
)

type TempTaskStorage struct {
	Storage sync.Map
}

func (s *TempTaskStorage) CreateTask(task domain.Task) error {
	if _, ok := s.Storage.Load(task.Id); ok {
		return fmt.Errorf("задача с таким id уже существует %s", task.Id.String())
	}
	s.Storage.Store(task.Id, task)
	return nil
}

func (s *TempTaskStorage) CreateTasks(tasks []domain.Task) error {
	for _, task := range tasks {
		if _, ok := s.Storage.Load(task.Id); ok {
			return fmt.Errorf("задача с таким id уже существует %s", task.Id.String())
		}
		s.Storage.Store(task.Id, task)
	}
	return nil
}

func (s *TempTaskStorage) GetTaskById(id uuid.UUID) (domain.Task, error) {
	if val, ok := s.Storage.Load(id); ok {
		return val.(domain.Task), nil
	}
	return domain.Task{}, fmt.Errorf("задачи с таким id не существует: %s", id.String())
}

func (s *TempTaskStorage) GetTasksByDayAndUser(day time.Weekday, chatID int64) ([]domain.Task, error) {
	result := make([]domain.Task, 100)
	for _, value := range s.Storage.Range {
		task := value.(domain.Task)
		if task.Day == day && task.ChatId == chatID {
			result = append(result, task)
		}
	}
	return result, nil
}

func (s *TempTaskStorage) GetCurrnetTask(chatID int64) (domain.Task, error) {
	for _, value := range s.Storage.Range {
		task := value.(domain.Task)
		day := time.Now().Weekday()
		curTime := time.Now().Local()
		totalMinutes := curTime.Hour()*60 + curTime.Minute()
		if task.Day == day && task.ChatId == chatID && task.Start <= totalMinutes && task.End > totalMinutes {
			return task, nil
		}
	}
	return domain.Task{}, nil
}

func (s *TempTaskStorage) UpdateTaskById(id uuid.UUID, task domain.Task) error {
	if _, ok := s.Storage.Load(id); ok {
		s.Storage.Store(id, task)
		return nil
	}
	return fmt.Errorf("задачи с таким id не существует: %s", id.String())
}

func (s *TempTaskStorage) DeleteTaskById(id uuid.UUID) error {
	s.Storage.LoadAndDelete(id)
	return nil
}

func (s *TempTaskStorage) DeleteTasksByDay(day time.Weekday, chatID int64) error {
	for _, value := range s.Storage.Range {
		task := value.(domain.Task)
		day := time.Now().Weekday()
		if task.Day == day && task.ChatId == chatID {
			s.Storage.LoadAndDelete(task.Id)
		}
	}
	return nil
}
