package app

import (
    "fmt"
    "log"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func processUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {

    if update.Message != nil {
        chatID := update.Message.Chat.ID
        text := update.Message.Text

        if text == "/start" {
            return sendWelcomeMessage(bot, chatID)
        }
    }

    if update.CallbackQuery != nil {
        query := update.CallbackQuery
        chatID := query.Message.Chat.ID
        data := query.Data

        switch data {

        // --- Первый этап: Нажатие кнопки "start" на клавиатуре ---
        case "action_start":
            // Отправляем основное меню из 4 кнопок
            return showMainMenu(bot, chatID, query.Message.MessageID)

        // --- Основное меню: 4 кнопки ---
        case "action_current_schedule":
            // Тут логика получения текущего расписания из базы
            scheduleText := "Текущее расписание (заглушка)."
            return editMessage(bot, chatID, query.Message.MessageID, scheduleText, mainMenuKeyboard())

        case "action_current_task":
            // Логика получения «текущей задачи»
            taskText := "Текущая задача (заглушка)."
            return editMessage(bot, chatID, query.Message.MessageID, taskText, mainMenuKeyboard())

        case "action_other_day":
            // Отобразим дни недели + кнопку «Назад»
            text := "Выберите день недели:"
            return editMessage(bot, chatID, query.Message.MessageID, text, daysOfWeekKeyboard())

        case "action_change_schedule":
            // Заглушка для изменения расписания
            text := "Здесь будет функционал изменения расписания."
            return editMessage(bot, chatID, query.Message.MessageID, text, mainMenuKeyboard())

        // --- Раскладка «Выберите день недели» ---
        case "action_back_to_main":
            // Вернуться к основному меню
            return editMessage(bot, chatID, query.Message.MessageID, "Меню:", mainMenuKeyboard())

        case "action_monday", "action_tuesday", "action_wednesday", "action_thursday",
             "action_friday", "action_saturday", "action_sunday":
            // Допустим, возвращаем расписание на выбранный день + снова основные кнопки
            dayName := parseDayName(data) // "Понедельник", "Вторник" и т.д.
            text := fmt.Sprintf("Расписание на %s (заглушка).", dayName)
            return editMessage(bot, chatID, query.Message.MessageID, text, mainMenuKeyboard())

        default:
            // Неизвестная callbackData — можно залогировать или игнорировать
            log.Println("Неизвестная callback data:", data)
        }
    }

    return nil
}

func sendWelcomeMessage(bot *tgbotapi.BotAPI, chatID int64) error {
    mainMenu := mainMenuKeyboard()

    msg := tgbotapi.NewMessage(chatID, "Выбери нужный пункт меню")
    msg.ReplyMarkup = mainMenu

    _, err := bot.Send(msg)
    return err
}

// showMainMenu - редактирует текущее сообщение, заменяя на основной набор кнопок
func showMainMenu(bot *tgbotapi.BotAPI, chatID int64, messageID int) error {
    text := "Выбери действие:"
    return editMessage(bot, chatID, messageID, text, mainMenuKeyboard())
}

// mainMenuKeyboard - возвращает InlineKeyboardMarkup с 4 основными кнопками
func mainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
    return tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("текущее расписание", "action_current_schedule"),
            tgbotapi.NewInlineKeyboardButtonData("текущая задача", "action_current_task"),
        ),
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("расписание на другой день", "action_other_day"),
            tgbotapi.NewInlineKeyboardButtonData("изменить расписание", "action_change_schedule"),
        ),
    )
}

func daysOfWeekKeyboard() tgbotapi.InlineKeyboardMarkup {
    // Для простоты сделаем 2 строки по 4 кнопки: Пн, Вт, Ср, Чт / Пт, Сб, Вс, Назад
    return tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("Понедельник", "action_monday"),
            tgbotapi.NewInlineKeyboardButtonData("Вторник",     "action_tuesday"),
            tgbotapi.NewInlineKeyboardButtonData("Среда",       "action_wednesday"),
            tgbotapi.NewInlineKeyboardButtonData("Четверг",     "action_thursday"),
        ),
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("Пятница",    "action_friday"),
            tgbotapi.NewInlineKeyboardButtonData("Суббота",    "action_saturday"),
            tgbotapi.NewInlineKeyboardButtonData("Воскресенье","action_sunday"),
            tgbotapi.NewInlineKeyboardButtonData("Назад",      "action_back_to_main"),
        ),
    )
}

func editMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int, newText string, newMarkup tgbotapi.InlineKeyboardMarkup) error {
    editMsg := tgbotapi.NewEditMessageText(chatID, messageID, newText)
    // Сначала редактируем текст
    _, err := bot.Send(editMsg)
    if err != nil {
        return err
    }

    // Затем редактируем клавиатуру
    newReplyMarkup := tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, newMarkup)
    _, err = bot.Send(newReplyMarkup)
    return err
}

func parseDayName(action string) string {
    switch action {
    case "action_monday":
        return "Понедельник"
    case "action_tuesday":
        return "Вторник"
    case "action_wednesday":
        return "Среда"
    case "action_thursday":
        return "Четверг"
    case "action_friday":
        return "Пятница"
    case "action_saturday":
        return "Суббота"
    case "action_sunday":
        return "Воскресенье"
    }
    return "Неизвестный день"
}
