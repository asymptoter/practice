package models

type User struct {
	Email    string `json:"email" gorm:"email"`
	Password string `json:"password" gorm:"password"`
	Token    string `json:"token" gorm:"token"`
}
