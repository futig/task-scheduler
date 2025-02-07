package service

func processUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, workerID int) {
    if update.Message == nil {
        return
    }

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
}