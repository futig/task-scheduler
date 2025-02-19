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

	userState, ok, err := service.GetUserState(wCfg.Cache, wCfg.Storage, chatID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("не удалось получить состояние пользователя %d", chatID)
	}

	if userState.State == enums.UpdateDaySchedule && regexp.MustCompile(`^[1-9]\d{0,3}$`).MatchString(text) {
		return showUpdateTask(bot, chatID, text, uuid.Nil, userState)
	} else if userState.State == enums.UpdateDaySchedule {
		return CreateSchedule(bot, chatID, text, userState)
	} else if userState.State == enums.UpdateTaskDescription {
		return UpdateTaskDescription(bot, chatID, text, userState)
	} else if userState.State == enums.UpdateTaskTime {
		return UpdateTaskTime(bot, chatID, text, userState)
	} else if userState.State == enums.UpdateTaskReminds {
		return UpdateTaskReminds(bot, chatID, text, userState)
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
		weekday := parseActionDay(data)
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
	err = createMessage(bot, chatID, "Выбери нужный пункт меню:", keyboard)
	if err != nil {
		return err
	}

	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
		ChatID:         chatID,
		State:          enums.MainMenu,
		ChoosenWeekday: -1,
		ChoosenItem:    uuid.Nil,
		LastMessageId:  -1,
	})
}

func showMainMenu(bot *tgbotapi.BotAPI, chatID int64, messageID int) (err error) {
	defer e.WrapErrIfNotNil(err, "showMainMenu")

	userState, ok, err := service.GetUserState(wCfg.Cache, wCfg.Storage, chatID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("не удалось получить состояние пользователя %d", chatID)
	}

	keyboard := mainMenuKeyboard()
	text := "Выбери нужный пункт меню:"

	if userState.LastMessageId != -1 && userState.LastMessageId != messageID {
		err = createMessage(bot, chatID, text, keyboard)
		messageID = -1
	} else {
		err = editMessage(bot, chatID, messageID, text, keyboard)
	}
	if err != nil {
		return err
	}

	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
		ChatID:         chatID,
		State:          enums.MainMenu,
		ChoosenWeekday: -1,
		ChoosenItem:    uuid.Nil,
		LastMessageId:  messageID,
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
	case enums.UpdateDaySchedule:
		return showChangeSchedule(bot, chatID, messageID)
	case enums.DeleteSchedule, enums.DayScheduleUpdated, enums.UpdateTask, enums.TaskNotFound:
		return showDayAction(bot, chatID, messageID, userState.ChoosenWeekday)
	case enums.DeleteTask, enums.UpdateTaskDescription, enums.UpdateTaskTime, enums.UpdateTaskReminds:
		return showUpdateTask(bot, chatID, "", userState.ChoosenItem, userState)
	}

	return nil
}

func showCurrentSchedule(bot *tgbotapi.BotAPI, chatID int64, messageID int) (err error) {
	defer e.WrapErrIfNotNil(err, "showCurrentSchedule")

	userState, ok, err := service.GetUserState(wCfg.Cache, wCfg.Storage, chatID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("не удалось получить состояние пользователя %d", chatID)
	}

	keyboard := infoKeyboard()
	text, ok, err := service.GetCurrentSchedule(wCfg.Storage, chatID)

	if err != nil {
		return err
	}

	if !ok {
		text = "На сегодня нет расписания."
	}

	if userState.LastMessageId != -1 && userState.LastMessageId != messageID {
		err = createMessage(bot, chatID, text, keyboard)
		messageID = -1
	} else {
		err = editMessage(bot, chatID, messageID, text, keyboard)
	}
	if err != nil {
		return err
	}

	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
		ChatID:         chatID,
		State:          enums.ShowCurrentSchedule,
		ChoosenWeekday: -1,
		ChoosenItem:    uuid.Nil,
		LastMessageId:  messageID,
	})
}

func showCurrentTask(bot *tgbotapi.BotAPI, chatID int64, messageID int) (err error) {
	defer e.WrapErrIfNotNil(err, "showCurrentTask")

	userState, ok, err := service.GetUserState(wCfg.Cache, wCfg.Storage, chatID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("не удалось получить состояние пользователя %d", chatID)
	}

	keyboard := infoKeyboard()
	tasks, ok, err := service.GetCurrentTasks(wCfg.Storage, chatID)

	if err != nil {
		return err
	}
	text := "Сейчас нет активных задач."
	if ok {
		text = tasks
	}

	if userState.LastMessageId != -1 && userState.LastMessageId != messageID {
		err = createMessage(bot, chatID, text, keyboard)
		messageID = -1
	} else {
		err = editMessage(bot, chatID, messageID, text, keyboard)
	}
	if err != nil {
		return err
	}

	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
		ChatID:         chatID,
		State:          enums.ShowCurrentTask,
		ChoosenWeekday: -1,
		ChoosenItem:    uuid.Nil,
		LastMessageId:  messageID,
	})
}

func showOtherDaySchedule(bot *tgbotapi.BotAPI, chatID int64, messageID int) (err error) {
	defer e.WrapErrIfNotNil(err, "showOtherDaySchedule")

	userState, ok, err := service.GetUserState(wCfg.Cache, wCfg.Storage, chatID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("не удалось получить состояние пользователя %d", chatID)
	}

	keyboard := daysOfWeekKeyboard()
	text := "Выберите день недели:"
	if userState.LastMessageId != -1 && userState.LastMessageId != messageID {
		err = createMessage(bot, chatID, text, keyboard)
		messageID = -1
	} else {
		err = editMessage(bot, chatID, messageID, text, keyboard)
	}
	if err != nil {
		return err
	}

	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
		ChatID:         chatID,
		State:          enums.ChooseWeekdayGet,
		ChoosenWeekday: -1,
		ChoosenItem:    uuid.Nil,
		LastMessageId:  messageID,
	})
}

func showChangeSchedule(bot *tgbotapi.BotAPI, chatID int64, messageID int) (err error) {
	defer e.WrapErrIfNotNil(err, "showChangeSchedule")

	userState, ok, err := service.GetUserState(wCfg.Cache, wCfg.Storage, chatID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("не удалось получить состояние пользователя %d", chatID)
	}

	keyboard := daysOfWeekKeyboard()
	text := "Выберите день недели для изменения расписания:"
	if userState.LastMessageId != -1 && userState.LastMessageId != messageID {
		err = createMessage(bot, chatID, text, keyboard)
		messageID = -1
	} else {
		err = editMessage(bot, chatID, messageID, text, keyboard)
	}
	if err != nil {
		return err
	}

	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
		ChatID:         chatID,
		State:          enums.ChooseWeekdayUpdate,
		ChoosenWeekday: -1,
		ChoosenItem:    uuid.Nil,
		LastMessageId:  messageID,
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
		schedule = "На выбранный день расписания нет."
	}

	if userState.State == enums.ChooseWeekdayGet {
		keyboard := infoKeyboard()
		if userState.LastMessageId != -1 && userState.LastMessageId != messageID {
			err = createMessage(bot, chatID, schedule, keyboard)
			messageID = -1
		} else {
			err = editMessage(bot, chatID, messageID, schedule, keyboard)
		}
		if err != nil {
			return err
		}
		return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
			ChatID:         chatID,
			State:          enums.ShowOtherDaySchedule,
			ChoosenWeekday: -1,
			ChoosenItem:    uuid.Nil,
			LastMessageId:  messageID,
		})
	}

	changeItemText := "Чтобы изменить пункт, отправь его номер."
	createText := `Чтобы добавить расписание, введи данные по примеру:

	--------------------
	7:00-17:30, работа по проекту, [15, 22 , 1], {10}
	13:00-15:21, просмотр кино, [2, 5], {30, 10}
	--------------------

	[] - за сколько минут нужно предупредить о начале задачи (от 1 до 59).
	{} - за сколько минут нужно предупредить о конце задачи (от 1 до 59).
	`
	var text string
	if ok {
		text = fmt.Sprintf("%s\n\n%s\n\n%s", schedule, changeItemText, createText)
	} else {
		text = fmt.Sprintf("%s\n\n%s", schedule, createText)
	}
	keyboard := changeScheduleKeyboard()
	if userState.LastMessageId != -1 && userState.LastMessageId != messageID {
		err = createMessage(bot, chatID, text, keyboard)
		messageID = -1
	} else {
		err = editMessage(bot, chatID, messageID, text, keyboard)
	}
	if err != nil {
		return err
	}
	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
		ChatID:         chatID,
		State:          enums.UpdateDaySchedule,
		ChoosenWeekday: weekday,
		ChoosenItem:    uuid.Nil,
		LastMessageId:  messageID,
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

	if userState.LastMessageId != -1 && userState.LastMessageId != messageID {
		err = createMessage(bot, chatID, text, keyboard)
		messageID = -1
	} else {
		err = editMessage(bot, chatID, messageID, text, keyboard)
	}
	if err != nil {
		return err
	}

	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
		ChatID:         chatID,
		State:          enums.DeleteSchedule,
		ChoosenWeekday: userState.ChoosenWeekday,
		ChoosenItem:    uuid.Nil,
		LastMessageId:  messageID,
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

	err = service.DeleteTaskById(wCfg.Storage, userState.ChoosenItem)
	if err != nil {
		return err
	}

	if userState.LastMessageId != -1 && userState.LastMessageId != messageID {
		err = createMessage(bot, chatID, text, keyboard)
		messageID = -1
	} else {
		err = editMessage(bot, chatID, messageID, text, keyboard)
	}
	if err != nil {
		return err
	}

	userState.State = enums.DeleteTask
	userState.LastMessageId = messageID
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

	if userState.LastMessageId != -1 && userState.LastMessageId != messageID {
		err = createMessage(bot, chatID, text, keyboard)
		messageID = -1
	} else {
		err = editMessage(bot, chatID, messageID, text, keyboard)
	}
	if err != nil {
		return err
	}

	userState.State = enums.UpdateTaskDescription
	userState.LastMessageId = messageID
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
	if userState.LastMessageId != -1 && userState.LastMessageId != messageID {
		err = createMessage(bot, chatID, text, keyboard)
		messageID = -1
	} else {
		err = editMessage(bot, chatID, messageID, text, keyboard)
	}
	if err != nil {
		return err
	}

	userState.State = enums.UpdateTaskTime
	userState.LastMessageId = messageID
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
	if userState.LastMessageId != -1 && userState.LastMessageId != messageID {
		err = createMessage(bot, chatID, text, keyboard)
		messageID = -1
	} else {
		err = editMessage(bot, chatID, messageID, text, keyboard)
	}
	if err != nil {
		return err
	}

	userState.State = enums.UpdateTaskReminds
	userState.LastMessageId = messageID
	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, userState)
}

func showUpdateTask(bot *tgbotapi.BotAPI, chatID int64, data string, taskId uuid.UUID, userState domain.User) (err error) {
	var task domain.Task
	var ok bool

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
		task, ok, err = service.GetTaskByPosition(wCfg.Storage, num, chatID, userState.ChoosenWeekday)
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
		return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
			ChatID:         chatID,
			State:          enums.UpdateTask,
			ChoosenWeekday: userState.ChoosenWeekday,
			ChoosenItem:    task.Id,
			LastMessageId:  -1,
		})
	}

	keyboard := infoKeyboard()
	err = createMessage(bot, chatID, "Указынный пункт не найден.", keyboard)
	if err != nil {
		return err
	}

	userState.State = enums.TaskNotFound
	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, userState)

}

func CreateSchedule(bot *tgbotapi.BotAPI, chatID int64, data string, userState domain.User) (err error) {
	userState, ok, err := service.GetUserState(wCfg.Cache, wCfg.Storage, chatID)
	if err != nil {
		return
	}
	if !ok {
		return fmt.Errorf("не удалось получить состояние пользователя %d", chatID)
	}
	if userState.State != enums.UpdateDaySchedule {
		return
	}

	err = service.CreateSchedule(wCfg.Storage, chatID, userState.ChoosenWeekday, data)
	if err != nil {
		return err
	}

	keyboard := infoKeyboard()
	err = createMessage(bot, chatID, "Расписание обновлено.", keyboard)
	if err != nil {
		return err
	}

	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
		ChatID:         chatID,
		State:          enums.DayScheduleUpdated,
		ChoosenWeekday: userState.ChoosenWeekday,
		ChoosenItem:    userState.ChoosenItem,
		LastMessageId:  -1,
	})
}

func UpdateTaskDescription(bot *tgbotapi.BotAPI, chatID int64, data string, userState domain.User) (err error) {
	err = service.UpdateTaskDescriptionById(wCfg.Storage, userState.ChoosenItem, data)
	if err != nil {
		return err
	}

	keyboard := infoKeyboard()
	err = createMessage(bot, chatID, "Описание обновлено.", keyboard)
	if err != nil {
		return err
	}

	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
		ChatID:         chatID,
		State:          enums.DayScheduleUpdated,
		ChoosenWeekday: userState.ChoosenWeekday,
		ChoosenItem:    userState.ChoosenItem,
		LastMessageId:  -1,
	})
}

func UpdateTaskTime(bot *tgbotapi.BotAPI, chatID int64, data string, userState domain.User) (err error) {
	err = service.UpdateTaskTimeById(wCfg.Storage, userState.ChoosenItem, data)
	if err != nil {
		return err
	}

	keyboard := infoKeyboard()
	err = createMessage(bot, chatID, "Время обновлено.", keyboard)
	if err != nil {
		return err
	}

	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
		ChatID:         chatID,
		State:          enums.DayScheduleUpdated,
		ChoosenWeekday: userState.ChoosenWeekday,
		ChoosenItem:    userState.ChoosenItem,
		LastMessageId:  -1,
	})
}

func UpdateTaskReminds(bot *tgbotapi.BotAPI, chatID int64, data string, userState domain.User) (err error) {
	err = service.UpdateTaskRemindsById(wCfg.Storage, userState.ChoosenItem, data)
	if err != nil {
		return err
	}

	keyboard := infoKeyboard()
	err = createMessage(bot, chatID, "Напоминания обновлены.", keyboard)
	if err != nil {
		return err
	}

	return service.SetUserState(wCfg.Cache, wCfg.Storage, chatID, domain.User{
		ChatID:         chatID,
		State:          enums.DayScheduleUpdated,
		ChoosenWeekday: userState.ChoosenWeekday,
		ChoosenItem:    userState.ChoosenItem,
		LastMessageId:  -1,
	})
}

// EDIT MESSAGE

func createMessage(bot *tgbotapi.BotAPI, chatID int64, newText string, keyboard tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatID, newText)
	msg.ReplyMarkup = keyboard

	_, err := bot.Send(msg)
	return err
}

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
	case "action_day_monday":
		return time.Monday
	case "action_day_tuesday":
		return time.Tuesday
	case "action_day_wednesday":
		return time.Wednesday
	case "action_day_thursday":
		return time.Thursday
	case "action_day_friday":
		return time.Friday
	case "action_day_saturday":
		return time.Saturday
	case "action_day_sunday":
		return time.Sunday
	}
	return -1
}
