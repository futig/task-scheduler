package service

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/futig/task-scheduler/internal/domain"
	"github.com/futig/task-scheduler/internal/domain/enums"
	"github.com/futig/task-scheduler/internal/storage"
	"github.com/google/uuid"
)

var taskRegex = regexp.MustCompile(`^(\d{1,2}):(\d{2})-(\d{1,2}):(\d{2}),\s+([^,]+),\s+\[(\s*(?:[1-5]?\d\s*,\s*)*[1-5]?\d\s*)\],\s+\{(\s*(?:[1-5]?\d\s*,\s*)*[1-5]?\d\s*)\}$`)
var taskTimeRegex = regexp.MustCompile(`^(\d{1,2}):(\d{2})-(\d{1,2}):(\d{2})$`)
var taskRemindsRegex = regexp.MustCompile(`^\[(\s*(?:[1-5]?\d\s*,\s*)*[1-5]?\d\s*)?\],\s*\{(\s*(?:[1-5]?\d\s*,\s*)*[1-5]?\d\s*)?\}$`)

func GetCurrentSchedule(storage storage.Storage, chatID int64) (string, bool, error) {
	return GetScheduleByWeekday(storage, chatID, time.Now().Local().Weekday())
}

func GetScheduleByWeekday(storage storage.Storage, chatID int64, weekday time.Weekday) (string, bool, error) {
	tasks, err := storage.GetTasksByDayAndUser(weekday, chatID)
	if err != nil {
		return "", false, err
	}
	if len(tasks) == 0 {
		return "", false, nil
	}
	res := make([]string, 0, len(tasks))
	for i, val := range tasks {
		res = append(res, fmt.Sprintf("%d) %s", i+1, val.String()))
	}
	schedule := strings.Join(res, "\n")
	return schedule, true, nil
}

func GetCurrentTasks(storage storage.Storage, chatID int64) (string, bool, error) {
	tasks, err := storage.GetCurrnetTasks(chatID)
	if err != nil {
		return "", false, err
	}
	if len(tasks) == 0 {
		return "", false, nil
	}
	res := make([]string, 0, len(tasks))
	for i, val := range tasks {
		res = append(res, fmt.Sprintf("%d) %s", i+1, val.String()))
	}
	schedule := strings.Join(res, "\n")
	return schedule, true, nil
}

func GetTaskById(storage storage.Storage, id uuid.UUID) (domain.Task, bool, error) {
	return storage.GetTaskById(id)
}

func GetTaskByPosition(storage storage.Storage, id int, chatID int64, weekday time.Weekday) (domain.Task, bool, error) {
	return storage.GetTaskByPosition(id-1, weekday, chatID)
}

func CreateSchedule(storage storage.Storage, chatID int64, weekday time.Weekday, data string) error {
	lines := strings.Split(data, "\n")
	tasks := make([]domain.Task, 0, len(lines))
	reminders := make([]domain.Remind, 0, len(lines))

	for _, line := range lines {
		matches := taskRegex.FindStringSubmatch(line)
		if matches == nil {
			return fmt.Errorf("invalid task format")
		}

		startHour, _ := strconv.Atoi(matches[1])
		startMin, _ := strconv.Atoi(matches[2])
		endHour, _ := strconv.Atoi(matches[3])
		endMin, _ := strconv.Atoi(matches[4])

		// Вычисляем минуты с начала дня
		startMinutes := startHour*60 + startMin
		endMinutes := endHour*60 + endMin

		// Комментарий
		comment := strings.TrimSpace(matches[5])

		// Создаём Task
		task := domain.Task{
			Id:      uuid.New(),
			ChatId:  chatID,
			Weekday: weekday,
			Start:   startMinutes,
			End:     endMinutes,
			Comment: comment,
		}
		tasks = append(tasks, task)

		parseReminders(matches[6], task.Id, enums.Start, &reminders, startMinutes)
		parseReminders(matches[7], task.Id, enums.End, &reminders, endMinutes)
		reminders = append(reminders, domain.Remind{
			Id:     uuid.New(),
			TaskId: task.Id,
			Type:   enums.Start,
			Time:   startMinutes,
		})
		reminders = append(reminders, domain.Remind{
			Id:     uuid.New(),
			TaskId: task.Id,
			Type:   enums.End,
			Time:   endMinutes,
		})
	}
	err := storage.CreateTasks(tasks)
	if err != nil {
		return err
	}
	return storage.CreateReminds(reminders)
}

func parseReminders(remindStr string, taskID uuid.UUID, remindType enums.RemindType, reminders *[]domain.Remind, baseTime int) {
	remindStr = strings.TrimSpace(remindStr)

	if remindStr == "" {
		return
	}

	times := strings.Split(remindStr, ",")

	for _, t := range times {
		t = strings.TrimSpace(t)
		minutes, err := strconv.Atoi(t)
		if err != nil || minutes < 1 || minutes > 59 {
			continue
		}

		remindTime := baseTime - minutes
		if remindTime < 0 {
			remindTime = 0 // Исключаем отрицательное время
		}

		*reminders = append(*reminders, domain.Remind{
			Id:     uuid.New(),
			TaskId: taskID,
			Type:   remindType,
			Time:   remindTime,
		})
	}
}

func UpdateTaskDescriptionById(storage storage.Storage, id uuid.UUID, data string) error {
	task, ok, err := storage.GetTaskById(id)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("не удалось найти задачу с таким id %s", id.String())
	}

	task.Comment = data
	return storage.UpdateTaskById(id, task)
}

func UpdateTaskTimeById(storage storage.Storage, id uuid.UUID, data string) error {
	matches := taskTimeRegex.FindStringSubmatch(data)
	if matches == nil {
		return fmt.Errorf("invalid task format")
	}

	startHour, _ := strconv.Atoi(matches[1])
	startMin, _ := strconv.Atoi(matches[2])
	endHour, _ := strconv.Atoi(matches[3])
	endMin, _ := strconv.Atoi(matches[4])

	startMinutes := startHour*60 + startMin
	endMinutes := endHour*60 + endMin

	task, ok, err := storage.GetTaskById(id)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("не удалось найти задачу с таким id %s", id.String())
	}

	task.Start = startMinutes
	task.End = endMinutes
	return storage.UpdateTaskById(id, task)
}

func UpdateTaskRemindsById(storage storage.Storage, id uuid.UUID, data string) error {
	matches := taskRemindsRegex.FindStringSubmatch(data)
	if matches == nil {
		return fmt.Errorf("invalid task format")
	}
	reminders := []domain.Remind{}

	task, ok, err := storage.GetTaskById(id)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("не удалось найти задачу с таким id %s", id.String())
	}

	parseReminders(matches[1], id, enums.Start, &reminders, task.Start)
	parseReminders(matches[2], id, enums.End, &reminders, task.End)

	err = storage.DeleteRemindsByTaskId(id)
	if err != nil {
		return err
	}
	return storage.CreateReminds(reminders)
}

func DeleteSchedule(storage storage.Storage, chatID int64, weekday time.Weekday) error {
	return storage.DeleteTasksByDay(weekday, chatID)
}

func DeleteTaskById(storage storage.Storage, id uuid.UUID) error {
	return storage.DeleteTaskById(id)
}
