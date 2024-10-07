package news_model

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
