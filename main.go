package main

import (
	"fmt"
	"os"
	"strconv"

	app "tgbot/internal"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {
	botToken := os.Getenv("TOKEN")

	// Note: Please keep in mind that default logger may expose sensitive information,
	// use in development only
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		return
	}
	app.Migrate()
	// Call method getMe
	botUser, _ := bot.GetMe()
	fmt.Printf("Bot User: %+v\n", botUser)

	updates, _ := bot.UpdatesViaLongPolling(nil)
	defer bot.StopLongPolling()

	// Inline keyboard parameters
	inlineKeyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow( // Row 1
			tu.InlineKeyboardButton("Добавить"). // Column 1
								WithCallbackData("add"),
			tu.InlineKeyboardButton("Удалить"). // Column 2
								WithCallbackData("remove"),
		),
		tu.InlineKeyboardRow( // Row 2
			tu.InlineKeyboardButton("Список"). // Column 1
								WithCallbackData("list"),
		),
	)

	for update := range updates {
		if update.CallbackQuery != nil {
			fmt.Println("Received callback with data:", update.CallbackQuery.Data)
			if update.CallbackQuery.Data == "list" {
				UserId := update.CallbackQuery.From.ID // User ID
				text := app.ReadList(UserId)
				send, _ := bot.SendMessage(tu.Message(tu.ID(UserId), text))
				fmt.Printf("Sent Message: %v\n", send)
			}
			if update.CallbackQuery.Data == "add" {
				UserId := update.CallbackQuery.From.ID // User ID
				app.MoveState(UserId, 1)
				send, _ := bot.SendMessage(tu.Message(tu.ID(UserId), "Напишите название лекарства"))
				fmt.Printf("Sent Message: %v\n", send)
			}
			if update.CallbackQuery.Data == "remove" {
				UserId := update.CallbackQuery.From.ID // User ID
				list := app.ReadList(UserId)
				app.MoveState(UserId, 99)
				send, _ := bot.SendMessage(tu.Message(tu.ID(UserId), "Список лекарств. \n\n"+list+"\n\n Введите номер лекарства (слева в списке) для удаления."))
				fmt.Printf("Sent Message: %v\n", send)
			}
		}
		if update.Message != nil {
			// Retrieve chat ID
			chatID := update.Message.Chat.ID

			if update.Message.Text == "/start" || update.Message.Text == "/help" {
				app.CreateUser(chatID)
				sentStartMessage, _ := bot.SendMessage(
					tu.Message(
						tu.ID(chatID),
						"Привет! Я бот, который умеет напоминать о таблетках.",
					).WithReplyMarkup(inlineKeyboard),
				)
				fmt.Printf("Sent Message: %v\n", sentStartMessage)
			}
			user := app.GetUser(chatID)
			if user.State > 0 {
				if user.State == 1 {
					app.AddMedicine(chatID, update.Message.Text)
					app.MoveState(chatID, 2)
					sentAddMessage, _ := bot.SendMessage(
						tu.Message(
							tu.ID(chatID),
							"Напишите время в формате 09:00, когда нужно принимать лекарство",
						),
					)
					fmt.Printf("Sent Message: %v\n", sentAddMessage)
				} else if user.State == 2 {
					err := app.AddTime(chatID, update.Message.Text)
					if err != nil {
						sentAddMessage, _ := bot.SendMessage(
							tu.Message(
								tu.ID(chatID),
								"Ошибка записи времени: \n\n"+err.Error(),
							),
						)
						fmt.Printf("Sent Message: %v\n", sentAddMessage)
					} else {
						app.MoveState(chatID, 0)
						pillString := app.GetPill(chatID)
						sentAddMessage, _ := bot.SendMessage(
							tu.Message(
								tu.ID(chatID),
								"Лекарство записано. \n\n"+pillString,
							),
						)
						fmt.Printf("Sent Message: %v\n", sentAddMessage)
					}
				} else if user.State == 99 {
					id, err := strconv.ParseInt(update.Message.Text, 10, 64)
					if err != nil {
						fmt.Println("Ошибка преобразования строки в число", err)
						sentAddMessage, _ := bot.SendMessage(
							tu.Message(
								tu.ID(chatID),
								"Некорректный ввод номера лекарства",
							),
						)
						fmt.Printf("Sent Message: %v\n", sentAddMessage)
					} else {
						app.RemovePill(chatID, id)
						sentAddMessage, _ := bot.SendMessage(
							tu.Message(
								tu.ID(chatID),
								"Лекарство под номером "+update.Message.Text+" удалено",
							),
						)
						app.MoveState(chatID, 0)
						fmt.Printf("Sent Message: %v\n", sentAddMessage)
					}

				}
			}

		}
	}
}
