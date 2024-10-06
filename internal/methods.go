package internal

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func Migrate() {
	d, db := DB()
	defer db.Close()
	d.AutoMigrate(&User{})
	d.AutoMigrate(&PillsData{})
}

func ReadList(UserId int64) string {
	d, db := DB()
	defer db.Close()
	var pills []PillsData
	d.Where("user_id = ?", UserId).Find(&pills)

	// Создаем срез строк для хранения отформатированных строк
	var result []string

	if len(pills) == 0 {
		return "Пока у вас нет таблеток"
	}
	// Итерируем по каждому элементу в срезе pills
	for _, pill := range pills {
		// Форматируем строку для каждого элемента
		formatted := fmt.Sprintf("%d: %s. Время приема: %s.", pill.Id, pill.Title, pill.Time)
		// Добавляем отформатированную строку в срез
		result = append(result, formatted)
	}

	// Объединяем все строки с разделителем "\n\n"
	return strings.Join(result, "\n\n")
}

func CreateUser(UserId int64) {
	d, db := DB()
	defer db.Close()
	user := User{UserId: UserId, State: 0, CurrentPillId: 0}
	d.Create(&user)
}

func GetUser(userId int64) User {
	d, db := DB()
	defer db.Close()
	var user User
	d.Where("user_id = ?", userId).First(&user)
	return user

}

func MoveState(userId int64, state int64) {
	d, db := DB()
	defer db.Close()
	var user User
	d.Where("user_id = ?", userId).First(&user)
	user.State = state
	d.Save(&user)
}

func AddMedicine(userId int64, title string) {
	d, db := DB()
	defer db.Close()
	var user User
	d.Where("user_id = ?", userId).First(&user)
	pill := PillsData{Title: title, UserId: userId}
	d.Create(&pill)
	user.CurrentPillId = pill.Id
	d.Save(&user)
}

func AddTime(userId int64, time string) error {
	d, db := DB()
	defer db.Close()
	var user User
	var pill PillsData
	d.Where("user_id = ?", userId).First(&user)
	d.Where("id = ?", user.CurrentPillId).First(&pill)
	if !isValidTimeString(time) {
		return fmt.Errorf("неверный формат времени, пример: 08:59")
	}
	pill.Time = time
	d.Save(&pill)
	return nil
}

func GetPill(userId int64) string {
	d, db := DB()
	defer db.Close()
	var pill PillsData
	var user User
	d.Where("user_id = ?", userId).First(&user)
	d.Where("id = ?", user.CurrentPillId).First(&pill)
	formatted := fmt.Sprintf("%d: %s. Время приема: %s.", pill.Id, pill.Title, pill.Time)
	return formatted

}

func isValidTimeString(s string) bool {
	// Регулярное выражение для проверки строки
	re := regexp.MustCompile(`^\d{2}:\d{2}$`)

	// Проверка длины строки
	if len(s) > 5 || len(s) < 5 {
		return false
	}
	//Проверка строки по регулярному выражению
	if !re.MatchString(s) {
		return false
	}

	// Извлечение числе из строки
	var hours, minuts int
	_, err := fmt.Sscanf(s, "%2d:%2d", &hours, &minuts)
	if err != nil {
		return false
	}
	// Проверка диапазона времени
	if hours < 0 || hours > 23 || minuts < 0 || minuts > 59 {
		return false
	}
	// Проверка строки по регулярному выражению
	return true
}

func RemovePill(userId int64, id int64) {
	d, db := DB()
	defer db.Close()
	var pill PillsData
	d.Where("id = ?", id).First(&pill)
	d.Delete(&pill)
}

func IsCurrentTime(hours, minutes int) bool {
	now := time.Now()

	return now.Hour() == hours && now.Minute() == minutes
}

func SendPillNotification(bot *telego.Bot) {
	d, db := DB()
	defer db.Close()
	var pills []PillsData
	d.Find(&pills)
	for _, pill := range pills {
		// Извлечение числе из строки
		var hours, minuts int
		_, err := fmt.Sscanf(pill.Time, "%2d:%2d", &hours, &minuts)
		if err != nil {
			return
		}
		if IsCurrentTime(hours, minuts) {
			send, _ := bot.SendMessage(tu.Message(tu.ID(pill.UserId), "Время принять лекарство: "+pill.Title))
			fmt.Printf("Sent Message: %v\n", send)
		}
	}
}
