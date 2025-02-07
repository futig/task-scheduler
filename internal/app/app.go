package app

import (
	"fmt"
	"log"
	"sync"

	config "github.com/futig/task-scheduler/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	cfg            = config.NewAppConfig()
	tasksCh        = make(chan tgbotapi.Update, cfg.QueueSize)
	stopWorkerCh   = make(chan struct{})
	currentWorkers = cfg.MinWorkers
	mu             sync.Mutex
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

	for i := 1; i <= currentWorkers; i++ {
		go worker(bot, i)
	}

	go manager(bot)

	log.Printf("Бот: %s запущен с %d воркерами (мин=%d, макс=%d)",
		bot.Self.UserName, currentWorkers, cfg.MinWorkers, cfg.MaxWorkers)

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
