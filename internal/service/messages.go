package service

import (
	"github.com/futig/task-scheduler/pkg/e"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TODO: вынести метод в app, сделать сервисы независимыми от телеграма
func SendMessage(bot *tgbotapi.BotAPI, chatID int64, message string) (err error) {	
	interimMsg := tgbotapi.NewMessage(chatID, message)
	_, err = bot.Send(interimMsg)
	err = e.WrapErrIfNotNil(err, "could not send message")
	return
}
