package service

// // handleMessage обрабатывает входящие текстовые сообщения в отдельной горутине
// func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
//     chatID := message.Chat.ID
//     text := message.Text

//     switch text {
//     case "/start":
//         // При старте можем отправить приветственное сообщение + кнопки
//         sendWelcome(bot, chatID)
//     default:
//         send(bot, chatID, "Неизвестная команда. Попробуй /start.")
//     }
// }

// // handleCallbackQuery обрабатывает нажатие кнопок
// func handleCallbackQuery(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) {
//     // Пример: при нажатии кнопки мы можем вернуть ответ пользователю
//     switch query.Data {
//     case "action_1":
//         answer := tgbotapi.NewCallback(query.ID, "Нажата кнопка 1!")
//         bot.Request(answer)

//     case "action_2":
//         answer := tgbotapi.NewCallback(query.ID, "Нажата кнопка 2!")
//         bot.Request(answer)
//     }

//     // Можно также изменять/обновлять сообщение, к которому прикреплена клавиатура
//     editMsg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID,
//         fmt.Sprintf("Вы нажали кнопку с данными: %s", query.Data))
//     bot.Send(editMsg)
// }

// // sendWelcome отправляет приветственное сообщение с «красивыми» кнопками
// func sendWelcome(bot *tgbotapi.BotAPI, chatID int64) {
//     text := "Привет! Нажми одну из кнопок ниже:"
//     // Создаём inline-клавиатуру
//     inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
//         tgbotapi.NewInlineKeyboardRow(
//             tgbotapi.NewInlineKeyboardButtonData("Кнопка 1", "action_1"),
//             tgbotapi.NewInlineKeyboardButtonData("Кнопка 2", "action_2"),
//         ),
//     )

//     msg := tgbotapi.NewMessage(chatID, text)
//     msg.ReplyMarkup = inlineKeyboard

//     bot.Send(msg)
// }