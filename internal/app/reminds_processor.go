package app

import (
	"fmt"
	"time"

	"github.com/futig/task-scheduler/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func processRemind(bot *tgbotapi.BotAPI, task domain.TaskRemind, workerId int) error {
	timeToTask := int(task.Start.Sub(time.Now()).Minutes())
	var message string
	if timeToTask == 0 {
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
		message = fmt.Sprintf("Через %d минут %s задача:\n%s", timeToTask, verb, task.Comment)
	}
	err := sendMessage(bot, task.ChatId, message)
	return err
}
