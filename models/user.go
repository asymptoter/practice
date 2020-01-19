package models

type User struct {
	ID               string `json:"ID" gorm:"id"`
	Email            string `json:"email" gorm:"email"`
	Password         string `json:"password" gorm:"password"`
	Token            string `json:"token" gorm:'token'`
	ActivationNumber int    `json:"activationNumber" gorm:"activation_number"`
}
