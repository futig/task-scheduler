package app

import (
	"context"
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func worker(ctx context.Context, bot *tgbotapi.BotAPI, workerID int, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("[Worker #%d] запущен\n", workerID)

	for {
		select {
		case update, ok := <-wCfg.UpdatesCh:
			if !ok {
				log.Printf("[Worker #%d] завершение (канал закрыт)\n", workerID)
				return
			}
			var chatID int64
			if update.Message != nil {
				chatID = update.Message.Chat.ID
			} else if update.CallbackQuery != nil {
				chatID = update.CallbackQuery.From.ID
			} else {
				continue
			}

			queueLen := len(wCfg.UpdatesCh) + len(wCfg.RemindsCh)
			if queueLen > cfg.BusyThreshold {
				_ = sendMessage(bot, chatID, "Запрос обрабатывается. Пожалуйста, подождите результата.")
			}
			
			processUpdate(bot, update)
			
		case task, ok := <-wCfg.RemindsCh:
			if !ok {
				log.Printf("[Worker #%d] завершение (канал закрыт)\n", workerID)
				return
			}
			processRemind(bot, task)

		case <-wCfg.StopWorkerCh:
			log.Printf("[Worker #%d] завершение (stopWorkerCh)\n", workerID)
			return
		case <-ctx.Done():
			log.Printf("[Worker #%d] завершение (stop signal)\n", workerID)
			return
		}
	}
}

func workersManager(ctx context.Context, bot *tgbotapi.BotAPI, wg *sync.WaitGroup) {
	defer wg.Done()

	currentWorkers := cfg.MinWorkers
	timerCh := time.After(cfg.WorkersCheckInterval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-timerCh:
			time.Sleep(2 * time.Second)

			wCfg.Mu.Lock()
			queueLen := len(wCfg.UpdatesCh) + len(wCfg.RemindsCh)

			if queueLen >= cfg.ScaleUpThreshold && currentWorkers < cfg.MaxWorkers {
				newWorkerCount := currentWorkers + 1
				wg.Add(1)
				go worker(ctx, bot, newWorkerCount, wg)
				currentWorkers = newWorkerCount
			}

			if queueLen <= cfg.ScaleDownThreshold && currentWorkers > cfg.MinWorkers {
				wCfg.StopWorkerCh <- struct{}{}
				currentWorkers--
			}

			wCfg.Mu.Unlock()
			timerCh = time.After(cfg.WorkersCheckInterval)
		}
	}
}
