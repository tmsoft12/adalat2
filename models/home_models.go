package model

type Views struct {
	ID     int `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID int `json:"user_id"`
	NewsID int `json:"news_id"`
}
type ViewsMedia struct {
	ID      int `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID  int `json:"user_id"`
	MediaID int `json:"media_id"`
}

type NewsSchema struct {
	ID             int    `json:"id" gorm:"primaryKey"`
	Image          string `json:"image"`
	TM_description string `json:"tm_description"`
	TM_title       string `json:"tm_title"`
	EN_title       string `json:"en_title"`
	RU_title       string `json:"ru_title"`
	EN_description string `json:"en_description"`
	RU_description string `json:"ru_description"`
	View           int    `json:"view"`
	Date           string `json:"date"`
}

type MediaSchema struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Cover    string `json:"cover"`
	Video    string `json:"video"`
	TM_title string `json:"tm_title"`
	EN_title string `json:"en_title"`
	RU_title string `json:"ru_title"`
	Date     string `json:"date"`
	View     int    `json:"view"`
}

type BannerSchema struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Image    string `json:"image"`
	Link     string `json:"link"`
	IsActive bool   `json:"is_active"`
}

type EmployerSchema struct {
	ID    int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name  string `json:"name"`
	Major string `json:"major"`
	Image string `json:"image"`
}
