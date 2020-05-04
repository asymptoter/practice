package models

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID `json:"ID" db:"id"`
	Token        uuid.UUID `json:"Token" db:"token"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"Email" db:"email"`
	Password     string    `json:"Password" db:"password"`
	RegisterDate int64     `json:"RegisterDate" db:"register_date"`
}
