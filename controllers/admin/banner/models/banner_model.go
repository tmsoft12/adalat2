package models

type BannerSchema struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Image    string `json:"image"`
	Link     string `json:"link"`
	IsActive bool   `json:"is_active"`
}
