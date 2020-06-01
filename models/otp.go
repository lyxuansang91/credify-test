package models

import (
	"time"
)

// OTP represent the OTP model
type OTP struct {
	ID        int64     `json:"id"`
	Ticket    string    `json:"ticket" validate:"required"`
	OTP       string    `json:"otp" validate:"required"`
	Phone     string    `json:"phone" validate:"required"`
	Verify    bool      `json:"verify"`
	ExpiredAt time.Time `json:"expired_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
