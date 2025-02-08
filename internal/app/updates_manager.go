package app

import (
	"context"
	"sync"

	"github.com/futig/task-scheduler/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func updatesManager(ctx context.Context, updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI, wg *sync.WaitGroup) {
	for {
		select {
		case update := <-updates:
			if update.Message != nil {
				chatID := update.Message.Chat.ID
				queueLen := len(wCfg.UpdatesCh) + len(wCfg.RemindsCh)

				if queueLen > cfg.BusyThreshold {
					_ = service.SendMessage(bot, chatID, "Запрос обрабатывается. Пожалуйста, подождите результата.")
				}
				
				wCfg.UpdatesCh <- update
			}
		case <-ctx.Done():
			wg.Done()
			return
		default:
			continue
		}
	}

}
