package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lyxuansang91/credify-test/models"
	"github.com/lyxuansang91/credify-test/user"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type mysqlUserRepository struct {
	Conn *sqlx.DB
}

// NewMysqlUserRepository will create an object that represent the user.Repository interface
func NewMysqlUserRepository(Conn *sqlx.DB) user.Repository {
	return &mysqlUserRepository{Conn}
}

func (m *mysqlUserRepository) GetByPhone(ctx context.Context, phone string) (res *models.User, err error) {
	query := `SELECT * FROM users WHERE active = true and phone = ?`
	results := []*models.User{}
	err = m.Conn.Select(&results, query, phone)
	if err != nil {
		return nil, err
	}

	if len(results) > 0 {
		res = results[0]
	} else {
		return nil, models.ErrNotFound
	}

	return
}

func (m *mysqlUserRepository) UpdateUser(ctx context.Context, u *models.User) error {
	qUpdateUser := `UPDATE users set first_name = ?, last_name = ?, email=?, country= ?, city=?, address_line=?, postal_code=?, province=?, updated_at = now() where phone = ?`
	result, err := m.Conn.Exec(qUpdateUser, u.FirstName, u.LastName, u.Email, u.Country, u.City, u.AddressLine, u.PostalCode, u.Province, u.Phone)
	if err != nil {
		return err
	}
	rowEffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowEffected == 0 {
		return errors.New("update information is not failed")
	}

	return nil
}

func (m *mysqlUserRepository) Store(ctx context.Context, u *models.User) error {
	query := `INSERT INTO users (phone, first_name, last_name, email, country, city, address_line, postal_code, created_at, updated_at, active) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.Active = true

	_, err := m.Conn.Exec(query, u.Phone, u.FirstName, u.LastName, u.Email, u.Country, u.City, u.AddressLine, u.PostalCode, u.CreatedAt, u.UpdatedAt, u.Active)
	return err
}
