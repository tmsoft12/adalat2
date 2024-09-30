package model

//	type MainSchema struct {
//		Banner   []BannerSchema   `json:"banner"`
//		News     []NewsSchema     `json:"news"`
//		Media    []MediaSchema    `json:"media"`
//		Employer []EmployerSchema `json:"employer"`
//	}
type Vi struct {
	ID     int `json:"id" gorm:"primaryKey;autoIncrement"` // Otomatik artan birincil anahtar
	UserID int `json:"user_id"`                            // Kullanıcı ID'si
	NewsID int `json:"news_id"`                            // Haber ID'si (bu ilişkilendirme için gerekli olabilir)
}

type NewsSchema struct {
	ID             int    `json:"id" gorm:"primaryKey"` // Birincil anahtar, genellikle int veya UUID olarak tanımlanır
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
	ID       string `json:"id"`
	Video    string `json:"video"`
	Title    string `json:"title"`
	EN_title string `json:"entitle"`
	RU_title string `json:"rutitle"`
	Date     string `json:"date"`
}

type BannerSchema struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Link        string `json:"link"`
}

type EmployerSchema struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Major string `json:"major"`
	Image string `json:"image"`
}
