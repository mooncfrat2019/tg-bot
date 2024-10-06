package internal

type User struct {
	Id            int64 `json:"id" gorm:"primaryKey"`
	UserId        int64 `json:"user_id"`
	CurrentPillId int64 `json:"current_pill_id"`
	State         int64 `json:"state"`
}

type PillsData struct {
	Id     int64  `json:"id" gorm:"primaryKey"`
	UserId int64  `json:"user_id"`
	Title  string `json:"title"`
	Time   string `json:"time"`
}
