package app

import (
	"github.com/futig/task-scheduler/pkg/error"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, message string) (err error) {	
	interimMsg := tgbotapi.NewMessage(chatID, message)
	_, err = bot.Send(interimMsg)
	err = e.WrapErrIfNotNil(err, "could not send message")
	return
}
