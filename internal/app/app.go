package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/futig/task-scheduler/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	cfg  = config.NewAppConfig()
	wCfg = config.NewWorkflowConfig()
)

func Run() {
	token := os.Getenv("BOT_TOKEN")
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

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup

	for i := 1; i <= cfg.MinWorkers; i++ {
		wg.Add(1)
		go worker(ctx, bot, i, &wg)
	}

	wg.Add(1)
	go workersManager(ctx, bot, &wg)

	wg.Add(1)
	go reminderManager(ctx, &wg)

	wg.Add(1)
	go updatesManager(ctx, updates, bot, &wg)

	log.Printf("Бот запущен")

	<-sigChan
	log.Println("Получен сигнал завершения, останавливаем бота...")

	cancel()

	wg.Wait()
	wCfg.CloseCh()

	log.Println("Бот остановлен")
}
