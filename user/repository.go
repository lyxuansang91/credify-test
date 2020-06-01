package user

import (
	"context"

	"github.com/lyxuansang91/credify-test/models"
)

// Repository represent the article's repository contract
type Repository interface {
	GetByPhone(ctx context.Context, phone string) (*models.User, error)
	Store(ctx context.Context, u *models.User) error
	UpdateUser(ctx context.Context, u *models.User) error
}
