package user

import (
	"context"

	"github.com/lyxuansang91/credify-test/models"
)

// Usecase represent the user's usecases
type Usecase interface {
	GetByPhone(ctx context.Context, phone string) (*models.User, error)
	Store(context.Context, *models.User) error
	VerifyOTP(ctx context.Context, phone string, otp string, ticket string) error
	CreateOTP(ctx context.Context, phone string) (*models.OTP, error)
}
