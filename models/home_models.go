package model

type MainSchema struct {
	Banner   []BannerSchema   `json:"banner"`
	News     []NewsSchema     `json:"news"`
	Media    []MediaSchema    `json:"media"`
	Employer []EmployerSchema `json:"employer"`
}

type NewsSchema struct {
	ID          string `json:"id" gorm:"primaryKeys"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Title       string `json:"title"`
	Date        string `json:"date"`
}

type MediaSchema struct {
	ID    string `json:"id"`
	Video string `json:"video"`
	Title string `json:"title"`
	Date  string `json:"date"`
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
