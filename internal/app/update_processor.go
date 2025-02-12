package app

import (
	"log"
	"regexp"
	"time"

	"github.com/futig/task-scheduler/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

func processUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {

	if update.Message != nil {
		return processMessage(bot, update.Message)
	}

	if update.CallbackQuery != nil {
		return processCallbackQuery(bot, update.CallbackQuery)
	}

	return nil
}

func processMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	chatID := message.Chat.ID
	text := message.Text

	if text == "/start" {
		return sendWelcomeMessage(bot, chatID)
	}

	return nil
}

func processCallbackQuery(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	chatID := query.Message.Chat.ID
	data := query.Data

	switch {

	case data == "action_start_menu":
		return showMainMenu(bot, chatID, query.Message.MessageID)

	case data == "action_current_schedule":
		return showCurrentSchedule(bot, chatID, query.Message.MessageID)

	case data == "action_current_task":
		return showCurrentTask(bot, chatID, query.Message.MessageID)

	case data == "action_other_day":
		return showOtherDaySchedule(bot, chatID, query.Message.MessageID)

	case data == "action_change_schedule":
		return showChangeSchedule(bot, chatID, query.Message.MessageID)

	case regexp.MustCompile(`^action_day`).MatchString(data):
		return showDayAction(bot, chatID, query.Message.MessageID, data)

	case regexp.MustCompile(`^action_delete_schedule`).MatchString(data):
		return showDeleteSchedule(bot, chatID, query.Message.MessageID, data)
		// return editMessage(bot, chatID, query.Message.MessageID, text, mainMenuKeyboard())

	case regexp.MustCompile(`^action_create_schedule`).MatchString(data):
		return showCreateSchedule(bot, chatID, query.Message.MessageID, data)
		// return editMessage(bot, chatID, query.Message.MessageID, text, mainMenuKeyboard())

	case data == "action_delete_schedule_item":
		return showDeleteScheduleItem(bot, chatID, query.Message.MessageID, data)
		// return editMessage(bot, chatID, query.Message.MessageID, text, mainMenuKeyboard())

	case data == "action_update_schedule_comment":
		return showUpdateCommentItem(bot, chatID, query.Message.MessageID, data)

	case data == "action_update_schedule_time":
		return showUpdateScheduleTime(bot, chatID, query.Message.MessageID, data)
		// return editMessage(bot, chatID, query.Message.MessageID, text, mainMenuKeyboard())

	case data == "action_update_schedule_reminds":
		return showUpdateScheduleReminds(bot, chatID, query.Message.MessageID, data)
		// return editMessage(bot, chatID, query.Message.MessageID, text, mainMenuKeyboard())

	default:
		log.Println("Неизвестная callback data:", data)
	}

	return nil
}

// KEBOARDS

func mainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Текущее расписание", "action_current_schedule"),
			tgbotapi.NewInlineKeyboardButtonData("Текущая задача", "action_current_task"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Расписание на другой день", "action_other_day"),
			tgbotapi.NewInlineKeyboardButtonData("Изменить расписание", "action_change_schedule"),
		),
	)
}

func infoKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "action_start_menu"),
		),
	)
}

func daysOfWeekKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Понедельник", "action_monday"),
			tgbotapi.NewInlineKeyboardButtonData("Вторник", "action_tuesday"),
			tgbotapi.NewInlineKeyboardButtonData("Среда", "action_wednesday"),
			tgbotapi.NewInlineKeyboardButtonData("Четверг", "action_thursday"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Пятница", "action_friday"),
			tgbotapi.NewInlineKeyboardButtonData("Суббота", "action_saturday"),
			tgbotapi.NewInlineKeyboardButtonData("Воскресенье", "action_sunday"),
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "action_start_menu"),
		),
	)
}

func createScheduleKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить расписание", "action_create_schedule"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "action_start_menu"),
		),
	)
}

func changeScheduleKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить расписание", "action_delete_schedule"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "action_start_menu"),
		),
	)
}

func changeScheduleItemKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить пункт", "action_delete_schedule_item"),
			tgbotapi.NewInlineKeyboardButtonData("Изменить описание", "action_update_schedule_comment"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить время", "action_update_schedule_time"),
			tgbotapi.NewInlineKeyboardButtonData("Изменить напоминания", "action_update_schedule_reminds"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "action_start_menu"),
		),
	)
}

// PRCESSORS

func sendWelcomeMessage(bot *tgbotapi.BotAPI, chatID int64) error {
	keyboard := mainMenuKeyboard()

	msg := tgbotapi.NewMessage(chatID, "Выбери нужный пункт меню:")
	msg.ReplyMarkup = keyboard

	_, err := bot.Send(msg)
	return err
}

func showMainMenu(bot *tgbotapi.BotAPI, chatID int64, messageID int) error {
	keyboard := mainMenuKeyboard()
	text := "Выбери нужный пункт меню:"
	return editMessageText(bot, chatID, messageID, text, keyboard)
}

func showCurrentSchedule(bot *tgbotapi.BotAPI, chatID int64, messageID int) error {
	keyboard := infoKeyboard()
	text, err := service.GetCurrentSchedule(chatID)

	if err != nil {
		return err
	}

	return editMessageText(bot, chatID, messageID, text, keyboard)
}

func showCurrentTask(bot *tgbotapi.BotAPI, chatID int64, messageID int) error {
	keyboard := infoKeyboard()
	text, err := service.GetCurrentTask(chatID)

	if err != nil {
		return err
	}

	return editMessageText(bot, chatID, messageID, text, keyboard)
}

func showOtherDaySchedule(bot *tgbotapi.BotAPI, chatID int64, messageID int) error {
	keyboard := daysOfWeekKeyboard()
	text := "Выберите день недели:"
	return editMessage(bot, chatID, messageID, text, keyboard)
}

func showChangeSchedule(bot *tgbotapi.BotAPI, chatID int64, messageID int) error {
	keyboard := daysOfWeekKeyboard()
	text := "Выберите день недели для изменения расписания:"
	return editMessage(bot, chatID, messageID, text, keyboard)
}

func showDayAction(bot *tgbotapi.BotAPI, chatID int64, messageID int, data string) error {

	keyboard := daysOfWeekKeyboard()
	text := "Выберите день недели для изменения расписания:"
	return editMessage(bot, chatID, messageID, text, keyboard)
}

func showCreateSchedule(bot *tgbotapi.BotAPI, chatID int64, messageID int, data string) error {
	keyboard := createScheduleKeyboard()
	text := `Чтобы добавить расписание, введи данные по примеру:
	
	7:00-17:30, работа по проекту, [15, 22 , 1], {10}
	13:00-15:21, просмотр кино, [2, 5], {30, 10}
	
	[] - за сколько минут нужно предупредить о начале задачи (от 1 до 59)
	{} - за сколько минут нужно предупредить о конце задачи (от 1 до 59)
	`
	return editMessage(bot, chatID, messageID, text, keyboard)
}

func showDeleteSchedule(bot *tgbotapi.BotAPI, chatID int64, messageID int, data string) error {
	keyboard := infoKeyboard()
	text := "Расписание удалено"
	err := service.DeleteSchedule(chatID, weekday)

	if err != nil {
		return err
	}

	return editMessage(bot, chatID, messageID, text, keyboard)
}

func showDeleteScheduleItem(bot *tgbotapi.BotAPI, chatID int64, messageID int, data string) error {
	id, err := uuid.Parse(data[13:])
	if err != nil {
		return err
	}

	weekday := parseDayName(data[13:])
	keyboard := infoKeyboard()
	text := "Пункт удален из расписания"

	err = service.DeleteScheduleItem(chatID, weekday, id)
	if err != nil {
		return err
	}

	return editMessage(bot, chatID, messageID, text, keyboard)
}

func showUpdateCommentItem(bot *tgbotapi.BotAPI, chatID int64, messageID int, data string) error {
	keyboard := infoKeyboard()
	text := "Введи обновленный комментарий"
	return editMessage(bot, chatID, messageID, text, keyboard)
}

func showUpdateScheduleTime(bot *tgbotapi.BotAPI, chatID int64, messageID int, data string) error {
	keyboard := infoKeyboard()
	text := "Введи обновленное время"
	return editMessage(bot, chatID, messageID, text, keyboard)
}

func showUpdateScheduleReminds(bot *tgbotapi.BotAPI, chatID int64, messageID int, data string) error {
	keyboard := infoKeyboard()
	text := `Введи обновленные напоминания по примеру:
	
	[15, 22 , 1], {10, 2}
	
	[] - за сколько минут нужно предупредить о начале задачи (от 1 до 59)
	{} - за сколько минут нужно предупредить о конце задачи (от 1 до 59)
	`
	return editMessage(bot, chatID, messageID, text, keyboard)
}

// EDIT MESSAGE

func editMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int, newText string, newMarkup tgbotapi.InlineKeyboardMarkup) error {
	err := editMessageText(bot, chatID, messageID, newText)
	if err != nil {
		return err
	}

	return editKeyboard(bot, chatID, messageID, newMarkup)
}

func editKeyboard(bot *tgbotapi.BotAPI, chatID int64, messageID int, newMarkup tgbotapi.InlineKeyboardMarkup) error {
	newReplyMarkup := tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, newMarkup)
	_, err := bot.Send(newReplyMarkup)
	return err
}

func editMessageText(bot *tgbotapi.BotAPI, chatID int64, messageID int, newText string) error {
	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, newText)
	_, err := bot.Send(editMsg)
	return err
}

// UTILS

func parseDayName(action string) time.Weekday {
	switch action {
	case "action_monday":
		return time.Monday
	case "action_tuesday":
		return time.Tuesday
	case "action_wednesday":
		return time.Wednesday
	case "action_thursday":
		return time.Thursday
	case "action_friday":
		return time.Friday
	case "action_saturday":
		return time.Saturday
	case "action_sunday":
		return time.Sunday
	}
	return -1
}
