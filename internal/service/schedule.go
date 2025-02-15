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

var taskRegex = regexp.MustCompile(`^\d{1,2}:\d{2}-\d{1,2}:\d{2},\s+[^,]+,\s+\[\s*(?:[1-5]?\d\s*,\s*)*[1-5]?\d\s*\],\s+\{\s*(?:[1-5]?\d\s*,\s*)*[1-5]?\d\s*\}$`)

func GetCurrentSchedule(storage storage.Storage, chatID int64) (string, bool, error) {
	return GetScheduleByWeekday(storage, chatID, time.Now().Local().Weekday())
}

func GetScheduleByWeekday(storage storage.Storage, chatID int64, weekday time.Weekday) (string, bool, error) {
	tasks, err := storage.GetTasksByDayAndUser(weekday, chatID)
	if err != nil {
		return "", false, err
	}
	res := make([]string, len(tasks))
	for _, val := range tasks {
		res = append(res, val.String())
	}
	schedule := strings.Join(res, "\n")
	return schedule, true, nil
}

func GetCurrentTasks(storage storage.Storage, chatID int64) (string, bool, error) {
	tasks, err := storage.GetCurrnetTasks(chatID)
	if err != nil {
		return "", false, err
	}
	res := make([]string, len(tasks))
	for _, val := range tasks {
		res = append(res, val.String())
	}
	schedule := strings.Join(res, "\n")
	return schedule, true, nil
}

func GetTaskById(storage storage.Storage, id uuid.UUID) (domain.Task, bool, error) {
	return storage.GetTaskById(id)
}

func GetTaskByPosition(storage storage.Storage, id int, chatID int64, weekday time.Weekday) (domain.Task, bool, error) {
	return storage.GetTaskByPosition(id, weekday, chatID)

}

func CreateSchedule(storage storage.Storage, chatID int64, weekday time.Weekday, data string) error {
	lines := strings.Split(data, "\n")
	tasks := make([]domain.Task, len(lines))
	reminders := make([]domain.Remind, len(lines))

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

func DeleteSchedule(storage storage.Storage, chatID int64, weekday time.Weekday) error {
	return storage.DeleteTasksByDay(weekday, chatID)
}

func DeleteTaskById(storage storage.Storage, id uuid.UUID) error {
	return storage.DeleteTaskById(id)
}
