package app

import (
	"fmt"
	"log"

	config "github.com/futig/task-scheduler/internal/config"
	service "github.com/futig/task-scheduler/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	cfg          = config.NewAppConfig()
	tasksCh      = make(chan tgbotapi.Update, cfg.QueueSize)
	stopWorkerCh = make(chan struct{})
)

func Run(token string) {
	if token == "" {
		log.Fatal("wrong token")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(fmt.Errorf("could not create connection with your token: %w", err))
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for i := 1; i <= cfg.MinWorkers; i++ {
		go service.Worker(bot, i, tasksCh, stopWorkerCh)
	}

	go service.WorkersManager(bot, tasksCh, stopWorkerCh, cfg)

	log.Printf("Бот: %s запущен с %d воркерами (мин=%d, макс=%d)",
		bot.Self.UserName, cfg.MinWorkers, cfg.MinWorkers, cfg.MaxWorkers)

	for update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID

			queueLen := len(tasksCh)
			if queueLen > cfg.BusyThreshold {
				interimMsg := tgbotapi.NewMessage(chatID, "Запрос обрабатывается. Пожалуйста, подождите результата.")
				bot.Send(interimMsg)
			}
			tasksCh <- update
		}
	}

	close(tasksCh)
	close(stopWorkerCh)
}
