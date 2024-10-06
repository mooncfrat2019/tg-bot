package main

import (
	"fmt"
	"os"
	"time"

	app "tgbot/internal"

	"github.com/mymmrac/telego"
	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New(cron.WithSeconds())
	botToken := os.Getenv("TOKEN")

	// Note: Please keep in mind that default logger may expose sensitive information,
	// use in development only
	bot, errBot := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if errBot != nil {
		fmt.Println(errBot)
	}
	// Добавление задачи, которая должна выполняться каждую минуту

	_, err := c.AddFunc("@every 1m", func() {
		app.SendPillNotification(bot)
		fmt.Println("Задача every1s выполнена:", time.Now())
	})

	if err != nil {
		fmt.Println("Ошибка при добавлении задачи:", err)
		return
	}

	// Запуск планировщика задач
	c.Start()

	// Оставляем приложение работающим
	select {}
}
