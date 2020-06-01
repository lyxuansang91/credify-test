package models

import (
	"time"
)

// User represent the user model
type User struct {
	ID          int64     `json:"id"`
	Phone       string    `json:"phone" validate:"required"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	Email       string    `json:"email" validate:"email"`
	Country     string    `json:"country"`
	Province    string    `json:"province"`
	City        string    `json:"city"`
	AddressLine string    `json:"address_line" db:"address_line"`
	PostalCode  string    `json:"postal_code" db:"postal_code"`
	Active      bool      `json:"active"`
	UpdatedAt   time.Time `json:"updated_at" db:"created_at"`
	CreatedAt   time.Time `json:"created_at" db:"updated_at"`
}
