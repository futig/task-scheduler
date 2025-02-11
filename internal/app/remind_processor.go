package app

import (
	"fmt"

	"github.com/futig/task-scheduler/internal/domain"
	"github.com/futig/task-scheduler/internal/domain/enums"
	t "github.com/futig/task-scheduler/pkg/time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func processRemind(bot *tgbotapi.BotAPI, task domain.TaskRemind) error {
	var timeToEvent int
	if task.Type == enums.Start {
		timeToEvent = task.Start - t.CurrentTimeToMinutes()
	} else {
		timeToEvent = task.End - t.CurrentTimeToMinutes()
	}

	var message string
	if timeToEvent == 0 {
		verb, err := task.Type.ParseToPresentVerb()
		if err != nil {
			return err
		}
		message = fmt.Sprintf("Задача %s:\n%s", verb, task.Comment)
	} else {
		verb, err := task.Type.ParseToFutureVerb()
		if err != nil {
			return err
		}
		message = fmt.Sprintf("Через %d минут %s задача:\n%s", timeToEvent, verb, task.Comment)
	}
	err := sendMessage(bot, task.ChatId, message)
	return err
}
