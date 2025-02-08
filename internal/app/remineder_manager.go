package app

import (
	"context"
	"sync"
	"time"

	"github.com/futig/task-scheduler/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func reminderManager(ctx context.Context, bot *tgbotapi.BotAPI, wg *sync.WaitGroup) {
	defer wg.Done()

	timerCh := time.After(cfg.RemindsCheckInterval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-timerCh:
			service.CheckAndSendReminders(wCfg.Storage, bot)
		default:
			continue
		}
	}
}
