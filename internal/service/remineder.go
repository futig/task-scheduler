package service

import (
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ReminderManager(bot *tgbotapi.BotAPI, mu *sync.Mutex, tasksCh chan tgbotapi.Update, stopWorkerCh chan struct{}) {
	for {
		time.Sleep(10 * time.Minute)
	}
}