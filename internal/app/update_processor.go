package app

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/futig/task-scheduler/internal/domain"
	"github.com/futig/task-scheduler/internal/domain/enums"
	"github.com/futig/task-scheduler/internal/service"
	e "github.com/futig/task-scheduler/pkg/error"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

func processUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message != nil {
		err := processMessage(bot, update.Message)
		if err != nil {
			sendMessage(bot, update.Message.Chat.ID, "Что пошло не так, попробуйте позже")
			log.Print(err)
		}
	}

	if update.CallbackQuery != nil {
		err := processCallbackQuery(bot, update.CallbackQuery)
		if err != nil {
			sendMessage(bot, update.CallbackQuery.Message.Chat.ID, "Что пошло не так, попробуйте позже")
			log.Print(err)
		}
	}
}

func processMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	chatID := message.Chat.ID
	text := message.Text

	if text == "/start" {
		return sendWelcomeMessage(bot, chatID)
	}
	if regexp.MustCompile(`^[1-9]\d{0,3}`).MatchString(text) {
		return showUpdateTask(bot, chatID, text, uuid.Nil)
	}

	return nil
}

func processCallbackQuery(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	chatID := query.Message.Chat.ID
	data := query.Data

	switch {

	case data == "action_start_menu":
		return showMainMenu(bot, chatID, query.Message.MessageID)

	case data == "action_back":
		return showPreviousAction(bot, chatID, query.Message.MessageID)

	case data == "action_current_schedule":
		return showCurrentSchedule(bot, chatID, query.Message.MessageID)

	case data == "action_current_task":
		return showCurrentTask(bot, chatID, query.Message.MessageID)

	case data == "action_other_day_schedule":
		return showOtherDaySchedule(bot, chatID, query.Message.MessageID)

	case data == "action_change_schedule":
		return showChangeSchedule(bot, chatID, query.Message.MessageID)

	case regexp.MustCompile(`^action_day`).MatchString(data):
		weekday := parseActionDay(data[11:])
		return showDayAction(bot, chatID, query.Message.MessageID, weekday)

	case data == "action_delete_schedule":
		return showDeleteSchedule(bot, chatID, query.Message.MessageID)

	case data == "action_delete_schedule_item":
		return showDeleteTask(bot, chatID, query.Message.MessageID)

	case data == "action_update_schedule_comment":
		return showUpdateTaskComment(bot, chatID, query.Message.MessageID)

	case data == "action_update_schedule_time":
		return showUpdateTaskTime(bot, chatID, query.Message.MessageID)

	case data == "action_update_schedule_reminds":
		return showUpdateTaskReminds(bot, chatID, query.Message.MessageID)

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
			tgbotapi.NewInlineKeyboardButtonData("Расписание на другой день", "action_other_day_schedule"),
			tgbotapi.NewInlineKeyboardButtonData("Изменить расписание", "action_change_schedule"),
		),
	)
}

func infoKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "action_back"),
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "action_start_menu"),
		),
	)
}

func daysOfWeekKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Понедельник", "action_day_monday"),
			tgbotapi.NewInlineKeyboardButtonData("Вторник", "action_day_tuesday"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Среда", "action_day_wednesday"),
			tgbotapi.NewInlineKeyboardButtonData("Четверг", "action_day_thursday"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Пятница", "action_day_friday"),
			tgbotapi.NewInlineKeyboardButtonData("Суббота", "action_day_saturday"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Воскресенье", "action_day_sunday"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "action_back"),
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
			tgbotapi.NewInlineKeyboardButtonData("Назад", "action_back"),
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
			tgbotapi.NewInlineKeyboardButtonData("Назад", "action_back"),
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "action_start_menu"),
		),
	)
}

// PRCESSORS

func sendWelcomeMessage(bot *tgbotapi.BotAPI, chatID int64) (err error) {
	defer e.WrapErrIfNotNil(err, "sendWelcomeMessage")

	keyboard := mainMenuKeyboard()

	msg := tgbotapi.NewMessage(chatID, "Выбери нужный пункт меню:")
	msg.ReplyMarkup = keyboard

	_, err = bot.Send(msg)
	if err != nil {
		return err
	}

	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
		ChatID:         chatID,
		State:          enums.MainMenu,
		ChoosenWeekday: -1,
		ChoosenItem:    uuid.Nil,
	})
}

func showMainMenu(bot *tgbotapi.BotAPI, chatID int64, messageID int) (err error) {
	defer e.WrapErrIfNotNil(err, "showMainMenu")

	keyboard := mainMenuKeyboard()
	text := "Выбери нужный пункт меню:"
	err = editMessage(bot, chatID, messageID, text, keyboard)
	if err != nil {
		return err
	}

	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
		ChatID:         chatID,
		State:          enums.MainMenu,
		ChoosenWeekday: -1,
		ChoosenItem:    uuid.Nil,
	})
}

func showPreviousAction(bot *tgbotapi.BotAPI, chatID int64, messageID int) (err error) {
	defer e.WrapErrIfNotNil(err, "showMainMenu")

	userState, ok, err := service.GetUserState(wCfg.Cache, wCfg.Storage, chatID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("не удалось получить состояние пользователя %d", chatID)
	}

	switch userState.State {
	case enums.ShowCurrentSchedule, enums.ShowCurrentTask, enums.ChooseWeekdayGet, enums.ChooseWeekdayUpdate:
		return showMainMenu(bot, chatID, messageID)
	case enums.ShowOtherDaySchedule:
		return showOtherDaySchedule(bot, chatID, messageID)
	case enums.UpdatingDaySchedule:
		return showChangeSchedule(bot, chatID, messageID)
	case enums.DeleteSchedule:
		return showChangeSchedule(bot, chatID, messageID)
	case enums.UpdateTask:
		return showDayAction(bot, chatID, messageID, userState.ChoosenWeekday)
	case enums.DeleteTask, enums.UpdateTaskDescription, enums.UpdateTaskTime, enums.UpdateTaskReminds:
		return showUpdateTask(bot, chatID, "", userState.ChoosenItem)
	}

	return nil
}

func showCurrentSchedule(bot *tgbotapi.BotAPI, chatID int64, messageID int) (err error) {
	defer e.WrapErrIfNotNil(err, "showCurrentSchedule")

	keyboard := infoKeyboard()
	text, ok, err := service.GetCurrentSchedule(wCfg.Storage, chatID)

	if err != nil {
		return err
	}

	if !ok {
		text = "На сегодня нет расписания"
	}

	err = editMessage(bot, chatID, messageID, text, keyboard)
	if err != nil {
		return err
	}

	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
		ChatID:         chatID,
		State:          enums.ShowCurrentSchedule,
		ChoosenWeekday: -1,
		ChoosenItem:    uuid.Nil,
	})
}

func showCurrentTask(bot *tgbotapi.BotAPI, chatID int64, messageID int) (err error) {
	defer e.WrapErrIfNotNil(err, "showCurrentTask")

	keyboard := infoKeyboard()
	task, ok, err := service.GetCurrentTask(wCfg.Storage, chatID)

	if err != nil {
		return err
	}
	text := "Сейчас нет активных задач."
	if ok {
		text = task.String()
	}

	err = editMessage(bot, chatID, messageID, text, keyboard)
	if err != nil {
		return err
	}

	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
		ChatID:         chatID,
		State:          enums.ShowCurrentTask,
		ChoosenWeekday: -1,
		ChoosenItem:    uuid.Nil,
	})
}

func showOtherDaySchedule(bot *tgbotapi.BotAPI, chatID int64, messageID int) (err error) {
	defer e.WrapErrIfNotNil(err, "showOtherDaySchedule")

	keyboard := daysOfWeekKeyboard()
	text := "Выберите день недели:"
	err = editMessage(bot, chatID, messageID, text, keyboard)
	if err != nil {
		return err
	}

	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
		ChatID:         chatID,
		State:          enums.ChooseWeekdayGet,
		ChoosenWeekday: -1,
		ChoosenItem:    uuid.Nil,
	})
}

func showChangeSchedule(bot *tgbotapi.BotAPI, chatID int64, messageID int) (err error) {
	defer e.WrapErrIfNotNil(err, "showChangeSchedule")

	keyboard := daysOfWeekKeyboard()
	text := "Выберите день недели для изменения расписания:"
	err = editMessage(bot, chatID, messageID, text, keyboard)
	if err != nil {
		return err
	}

	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
		ChatID:         chatID,
		State:          enums.ChooseWeekdayUpdate,
		ChoosenWeekday: -1,
		ChoosenItem:    uuid.Nil,
	})
}

func showDayAction(bot *tgbotapi.BotAPI, chatID int64, messageID int, weekday time.Weekday) (err error) {
	defer e.WrapErrIfNotNil(err, fmt.Sprintf("showDayAction %s", weekday.String()))

	userState, ok, err := service.GetUserState(wCfg.Cache, wCfg.Storage, chatID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("не удалось получить состояние пользователя %d", chatID)
	}

	schedule, ok, err := service.GetScheduleByWeekday(wCfg.Storage, chatID, weekday)
	if err != nil {
		return err
	}
	if !ok {
		schedule = "На выбранный день расписания нет"
	}

	if userState.State == enums.ChooseWeekdayGet {
		keyboard := infoKeyboard()
		err = editMessage(bot, chatID, messageID, schedule, keyboard)
		if err != nil {
			return err
		}
		return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
			ChatID:         chatID,
			State:          enums.ShowOtherDaySchedule,
			ChoosenWeekday: -1,
			ChoosenItem:    uuid.Nil,
		})
	}

	changeItemText := "Чтобы изменить пункт, отправь его номер."
	createText := `Чтобы добавить расписание, введи данные по примеру:

	--------------------
	7:00-17:30, работа по проекту, [15, 22 , 1], {10}
	13:00-15:21, просмотр кино, [2, 5], {30, 10}
	--------------------

	[] - за сколько минут нужно предупредить о начале задачи (от 1 до 59)
	{} - за сколько минут нужно предупредить о конце задачи (от 1 до 59)
	`
	var text string
	if ok {
		text = fmt.Sprintf("%s\n\n%s", schedule, createText)
	} else {
		text = fmt.Sprintf("%s\n\n%s\n\n%s", schedule, changeItemText, createText)
	}
	keyboard := changeScheduleKeyboard()
	err = editMessage(bot, chatID, messageID, text, keyboard)
	if err != nil {
		return err
	}
	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
		ChatID:         chatID,
		State:          enums.UpdatingDaySchedule,
		ChoosenWeekday: weekday,
		ChoosenItem:    uuid.Nil,
	})
}

func showDeleteSchedule(bot *tgbotapi.BotAPI, chatID int64, messageID int) (err error) {
	keyboard := infoKeyboard()
	text := "Расписание удалено"

	userState, ok, err := service.GetUserState(wCfg.Cache, wCfg.Storage, chatID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("не удалось получить состояние пользователя %d", chatID)
	}

	err = service.DeleteSchedule(wCfg.Storage, chatID, userState.ChoosenWeekday)
	if err != nil {
		return err
	}

	err = editMessage(bot, chatID, messageID, text, keyboard)
	if err != nil {
		return err
	}

	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
		ChatID:         chatID,
		State:          enums.DeleteSchedule,
		ChoosenWeekday: userState.ChoosenWeekday,
		ChoosenItem:    uuid.Nil,
	})
}

func showDeleteTask(bot *tgbotapi.BotAPI, chatID int64, messageID int) (err error) {
	userState, ok, err := service.GetUserState(wCfg.Cache, wCfg.Storage, chatID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("не удалось получить состояние пользователя %d", chatID)
	}

	keyboard := infoKeyboard()
	text := "Пункт удален из расписания"

	err = service.DeleteScheduleItem(wCfg.Storage, chatID, userState.ChoosenWeekday, userState.ChoosenItem)
	if err != nil {
		return err
	}

	err = editMessage(bot, chatID, messageID, text, keyboard)
	if err != nil {
		return err
	}

	userState.State = enums.DeleteTask
	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, userState)
}

func showUpdateTaskComment(bot *tgbotapi.BotAPI, chatID int64, messageID int) (err error) {
	userState, ok, err := service.GetUserState(wCfg.Cache, wCfg.Storage, chatID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("не удалось получить состояние пользователя %d", chatID)
	}

	keyboard := infoKeyboard()
	text := "Введи обновленный комментарий"

	err = editMessage(bot, chatID, messageID, text, keyboard)
	if err != nil {
		return err
	}

	userState.State = enums.UpdateTaskDescription
	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, userState)
}

func showUpdateTaskTime(bot *tgbotapi.BotAPI, chatID int64, messageID int) (err error) {
	userState, ok, err := service.GetUserState(wCfg.Cache, wCfg.Storage, chatID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("не удалось получить состояние пользователя %d", chatID)
	}

	keyboard := infoKeyboard()
	text := "Введи обновленное время"
	err = editMessage(bot, chatID, messageID, text, keyboard)
	if err != nil {
		return err
	}

	userState.State = enums.UpdateTaskTime
	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, userState)
}

func showUpdateTaskReminds(bot *tgbotapi.BotAPI, chatID int64, messageID int) (err error) {
	userState, ok, err := service.GetUserState(wCfg.Cache, wCfg.Storage, chatID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("не удалось получить состояние пользователя %d", chatID)
	}

	keyboard := infoKeyboard()
	text := `Введи обновленные напоминания по примеру:
	
	[15, 22 , 1], {10, 2}
	
	[] - за сколько минут нужно предупредить о начале задачи (от 1 до 59)
	{} - за сколько минут нужно предупредить о конце задачи (от 1 до 59)
	`
	err = editMessage(bot, chatID, messageID, text, keyboard)
	if err != nil {
		return err
	}

	userState.State = enums.UpdateTaskReminds
	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, userState)
}

func showUpdateTask(bot *tgbotapi.BotAPI, chatID int64, data string, taskId uuid.UUID) (err error) {
	userState, ok, err := service.GetUserState(wCfg.Cache, wCfg.Storage, chatID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("не удалось получить состояние пользователя %d", chatID)
	}

	var task domain.Task

	if taskId != uuid.Nil {
		task, ok, err = service.GetTaskById(wCfg.Storage, taskId)
		if err != nil {
			return err
		}
	} else {
		num, err := strconv.Atoi(data)
		if err != nil {
			return err
		}
		task, ok, err = service.GetTaskByPosition(wCfg.Storage, num)
		if err != nil {
			return err
		}
	}

	if ok {
		keyboard := changeScheduleItemKeyboard()
		msg := tgbotapi.NewMessage(chatID, task.String())
		msg.ReplyMarkup = keyboard
		if err != nil {
			return err
		}

		_, err = bot.Send(msg)
		if err != nil {
			return err
		}
		taskId = task.Id
		userState.State = enums.UpdateTask
		userState.ChoosenItem = taskId
		return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, userState)
	}

	keyboard := infoKeyboard()
	msg := tgbotapi.NewMessage(chatID, "Указынный пункт не найден.")
	msg.ReplyMarkup = keyboard
	if err != nil {
		return err
	}

	_, err = bot.Send(msg)
	if err != nil {
		return err
	}

	userState.State = enums.TaskNotFound
	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, userState)

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

func parseActionDay(action string) time.Weekday {
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
