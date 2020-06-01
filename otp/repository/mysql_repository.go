package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lyxuansang91/credify-test/models"
	"github.com/lyxuansang91/credify-test/otp"
)

type mysqlOTPRepo struct {
	Conn *sqlx.DB
}

// NewMysqlOTPRepository will create an implementation of otp.Repository
func NewMysqlOTPRepository(db *sqlx.DB) otp.Repository {
	return &mysqlOTPRepo{
		Conn: db,
	}
}

func (m *mysqlOTPRepo) GetByID(ctx context.Context, id int64) (*models.OTP, error) {
	query := `SELECT id, ticket, phone, otp FROM otps WHERE id = ?`
	otp := &models.OTP{}
	err := m.Conn.Get(&otp, query, id)
	if err != nil {
		return nil, err
	}
	if otp == nil {
		return nil, errors.New("OTP not found")
	}
	return otp, nil
}

func (m *mysqlOTPRepo) GetByTicket(ctx context.Context, ticket string) (*models.OTP, error) {
	query := `SELECT id, ticket, phone, otp FROM otps WHERE ticket = ?`
	otps := []*models.OTP{}
	err := m.Conn.Select(&otps, query, ticket)
	if err != nil {
		return nil, err
	}
	if otps == nil {
		return nil, errors.New("OTP is not found")
	}
	otp := &models.OTP{}
	if len(otps) > 0 {
		otp = otps[0]
	}
	return otp, nil
}

func (m *mysqlOTPRepo) CreateOTP(ctx context.Context, otp *models.OTP) error {
	qInsertOTP := `INSERT INTO otps (ticket, phone, otp, expired_at, created_at) VALUES (?, ?, ?, ?, ?)`
	otp.CreatedAt = time.Now()
	otp.ExpiredAt = otp.CreatedAt.Add(5 * time.Minute) // expired 5 min
	_, err := m.Conn.Exec(qInsertOTP, otp.Ticket, otp.Phone, otp.OTP, otp.ExpiredAt, otp.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (m *mysqlOTPRepo) VerifyOTP(ctx context.Context, phone, otp, ticket string) error {
	foundOTP, err := m.GetByTicket(ctx, ticket)
	if err != nil {
		return err
	}
	if foundOTP == nil {
		return errors.New("OTP not found")
	}

	if foundOTP.ExpiredAt.After(time.Now()) {
		return errors.New("OTP is expired")
	}

	if foundOTP.Phone != phone || foundOTP.OTP != otp || foundOTP.Ticket != ticket {
		return errors.New("OTP verification failed")
	}
	return nil
}
