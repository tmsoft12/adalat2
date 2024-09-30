package model

type MainSchema struct {
	Banner   []BannerSchema   `json:"banner"`
	News     []NewsSchema     `json:"news"`
	Media    []MediaSchema    `json:"media"`
	Employer []EmployerSchema `json:"employer"`
}
type Views struct {
	ID     int    `json:"id"`
	UserID string `josn:"userid"`
}
type NewsSchema struct {
	ID             string `json:"id" gorm:"primaryKeys"`
	Image          string `json:"image"`
	Description    string `json:"description"`
	Title          string `json:"title"`
	Date           string `json:"date"`
	EN_title       string `json:"entitle"`
	RU_title       string `json:"rutitle"`
	EN_description string `json:"endescription"`
	RU_description string `json:"rudescription"`
	Count          string `jsin:"count"`
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
