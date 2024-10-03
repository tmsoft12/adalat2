package model

type Vi struct {
	ID     int `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID int `json:"user_id"`
	NewsID int `json:"news_id"`
}

type NewsSchema struct {
	ID             int    `json:"id" gorm:"primaryKey"`
	Image          string `json:"image"`
	Description    string `json:"description"`
	Title          string `json:"title"`
	Date           string `json:"date"`
	EN_title       string `json:"entitle"`
	RU_title       string `json:"rutitle"`
	EN_description string `json:"endescription"`
	RU_description string `json:"rudescription"`
	Count          int    `json:"count"`
}

type MediaSchema struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Video    string `json:"video"`
	Title    string `json:"title"`
	EN_title string `json:"entitle"`
	RU_title string `json:"rutitle"`
	Date     string `json:"date"`
}

type BannerSchema struct {
	ID    int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Image string `json:"image"`
	Link  string `json:"link"`
}

type EmployerSchema struct {
	ID    int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name  string `json:"name"`
	Major string `json:"major"`
	Image string `json:"image"`
}
