package employer_models

type EmployerSchema struct {
	ID      int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Major   string `json:"major"`
	Image   string `json:"image"`
}
