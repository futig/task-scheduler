package service

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


func ProcessUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, workerID int) error {
    chatID := update.Message.Chat.ID
    text := update.Message.Text

    // Допустим, имитируем "тяжёлую" операцию
    time.Sleep(1 * time.Second)

    respText := fmt.Sprintf("Ваш запрос обработан в Worker #%d. Текст: %s", workerID, text)
    msg := tgbotapi.NewMessage(chatID, respText)
    _, err := bot.Send(msg)
    if err != nil {
        log.Printf("[Worker #%d] ошибка отправки ответа: %v", workerID, err)
    }
    return nil
}