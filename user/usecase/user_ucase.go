package usecase

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/lyxuansang91/credify-test/models"
	"github.com/lyxuansang91/credify-test/otp"
	"github.com/lyxuansang91/credify-test/user"
)

type userUsecase struct {
	userRepo       user.Repository
	otpRepo        otp.Repository
	contextTimeout time.Duration
}

// NewUserUsecase will create new an userUsecase object representation of article.Usecase interface
func NewUserUsecase(u user.Repository, o otp.Repository, timeout time.Duration) user.Usecase {
	return &userUsecase{
		userRepo:       u,
		otpRepo:        o,
		contextTimeout: timeout,
	}
}

// Store ...
func (u *userUsecase) Store(c context.Context, m *models.User) error {

	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	newUser, err := u.userRepo.GetByPhone(ctx, m.Phone)
	if err != nil {
		return err
	}
	if newUser == nil {
		return errors.New("user is not verified")
	}

	if err = u.userRepo.UpdateUser(ctx, m); err != nil {
		return err
	}

	return nil
}

// GetByPhone ...
func (u *userUsecase) GetByPhone(c context.Context, phone string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	user, err := u.userRepo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("Not found")
	}
	return user, nil
}

// GetByPhone ...
func (u *userUsecase) VerifyOTP(c context.Context, phone string, otp string, ticket string) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.otpRepo.VerifyOTP(ctx, phone, otp, ticket)
}

func randStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (u *userUsecase) CreateOTP(c context.Context, phone string) (*models.OTP, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	otp := &models.OTP{
		Phone:  phone,
		OTP:    randStringRunes(6),
		Ticket: randStringRunes(10),
	}
	if err := u.otpRepo.CreateOTP(ctx, otp); err != nil {
		return nil, err
	}
	newOtp, erro := u.otpRepo.GetByTicket(ctx, otp.Ticket)
	if erro != nil {
		return nil, erro
	}
	if newOtp == nil {
		return nil, errors.New("OTP not found")
	}
	return newOtp, nil
}
