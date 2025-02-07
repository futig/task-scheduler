package service

import (
	"log"
	"sync"
	"time"

	"github.com/futig/task-scheduler/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Worker(bot *tgbotapi.BotAPI, workerID int, tasksCh chan tgbotapi.Update, stopWorkerCh chan struct{}) {
	log.Printf("[Worker #%d] запущен\n", workerID)
	for {
		select {
		case update, ok := <-tasksCh:
			if !ok {
				log.Printf("[Worker #%d] завершение (канал закрыт)\n", workerID)
				return
			}
			processUpdate(bot, update, workerID)
		case <-stopWorkerCh:
			log.Printf("[Worker #%d] завершение (stopWorkerCh)\n", workerID)
			return
		}
	}
}

func WorkersManager(bot *tgbotapi.BotAPI, tasksCh chan tgbotapi.Update, stopWorkerCh chan struct{}, cfg config.AppConfig) {
	currentWorkers := cfg.MinWorkers
	var mu *sync.Mutex

	for {
        time.Sleep(2 * time.Second)

        mu.Lock()
        queueLen := len(tasksCh)

		if queueLen >= cfg.ScaleUpThreshold && currentWorkers < cfg.MaxWorkers {
			newWorkerCount := currentWorkers + 1
			log.Printf("[Manager] Увеличиваем воркеров с %d до %d. Длина очереди=%d\n",
				currentWorkers, newWorkerCount, queueLen)
			go Worker(bot, newWorkerCount, tasksCh, stopWorkerCh)
			currentWorkers = newWorkerCount
		}

		if queueLen <= cfg.ScaleDownThreshold && currentWorkers > cfg.MinWorkers {
			log.Printf("[Manager] Уменьшаем воркеров с %d до %d. Длина очереди=%d\n",
				currentWorkers, currentWorkers-1, queueLen)
			stopWorkerCh <- struct{}{}
			currentWorkers--
		}

		mu.Unlock()
	}
}
