package models

type User struct {
	ID           string `json:"ID" db:"id"`
	Email        string `json:"Email" db:"email"`
	Password     string `db:"password"`
	Token        string `json:"Token" db:"token"`
	RegisterDate int64  `json:"RegisterDate" db:"register_date"`
}
