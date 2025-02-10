package tempstorage

import (
	"fmt"

	"github.com/futig/task-scheduler/internal/domain"
	t "github.com/futig/task-scheduler/pkg/time"

	"github.com/google/uuid"
)

func (s *TempStorageContext) CreateRemind(remind domain.Remind) error {
	if _, ok := s.Reminds.Load(remind.Id); ok {
		return fmt.Errorf("напоминание с таким id уже существует %s", remind.Id.String())
	}
	s.Reminds.Store(remind.Id, remind)
	return nil
}

func (s *TempStorageContext) CreateReminds(reminds []domain.Remind) error {
	for _, remind := range reminds {
		if _, ok := s.Reminds.Load(remind.Id); ok {
			return fmt.Errorf("напоминание с таким id уже существует %s", remind.Id.String())
		}
		s.Reminds.Store(remind.Id, remind)
	}
	return nil
}

func (s *TempStorageContext) GetRemindById(id uuid.UUID) (domain.TaskRemind, bool, error) {
	if remVal, ok := s.Reminds.Load(id); ok {
		remind := remVal.(domain.Remind)
		if taskVal, ok := s.Reminds.Load(remind.TaskId); ok {
			task := taskVal.(domain.Task)
			taskRem := domain.TaskRemind{
				ChatId:  task.ChatId,
				Start:   task.Start,
				End:     task.End,
				Comment: task.Comment,
				Type:    remind.Type,
			}
			return taskRem, true, nil
		}
	}
	return domain.TaskRemind{}, false, fmt.Errorf("напоминания с таким id не существует: %s", id.String())
}

func (s *TempStorageContext) GetRemindsForPeriod(period int) ([]domain.TaskRemind, error) {
	result := make([]domain.TaskRemind, 100)
	for _, remVal := range s.Reminds.Range {
		remind := remVal.(domain.Remind)
		curTime := t.CurrentTimeToMinutes()
		diff := remind.Time - curTime
		if diff < 0 && period <= diff {
			continue
		}
		if taskVal, ok := s.Reminds.Load(remind.TaskId); ok {
			task := taskVal.(domain.Task)
			taskRem := domain.TaskRemind{
				ChatId:  task.ChatId,
				Start:   task.Start,
				End:     task.End,
				Comment: task.Comment,
				Type:    remind.Type,
			}
			result = append(result, taskRem)
		}
	}
	return result, nil
}

func (s *TempStorageContext) DeleteRemindById(id uuid.UUID) error {
	s.Reminds.Delete(id)
	return nil
}

func (s *TempStorageContext) DeleteRemindsByTaskId(id uuid.UUID) error {
	for _, value := range s.Reminds.Range {
		task := value.(domain.Task)
		if task.Id == id {
			s.Reminds.Delete(task.Id)
		}
	}
	return nil
}
