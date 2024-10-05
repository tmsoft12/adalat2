package media_model

type MediaSchema struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Video    string `json:"video"`
	TM_title string `json:"tm_title"`
	EN_title string `json:"en_title"`
	RU_title string `json:"ru_title"`
	Date     string `json:"date"`
}
