package otp

import (
	"context"

	"github.com/lyxuansang91/credify-test/models"
)

// Repository represent the author's repository contract
type Repository interface {
	GetByID(ctx context.Context, id int64) (*models.OTP, error)
	GetByTicket(ctx context.Context, ticket string) (*models.OTP, error)
	CreateOTP(ctx context.Context, otp *models.OTP) error
	VerifyOTP(ctx context.Context, phone, otp, ticket string) error
}
