package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/lyxuansang91/credify-test/models"
	"github.com/lyxuansang91/credify-test/user"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// UserHandler  represent the httphandler for user
type UserHandler struct {
	UUsecase user.Usecase
}

// NewUserHandler will initialize the users/ resources endpoint
func NewUserHandler(e *echo.Echo, us user.Usecase) {
	handler := &UserHandler{
		UUsecase: us,
	}
	e.POST("/users", handler.Store)
	// e.GET("/users/user-by-phone", handler.GetByPhone)
	e.POST("/users/verify-otp", handler.VerifyOTP)
	e.POST("/users/create-otp", handler.CreateOTP)
}

// CreateOTP will create otp when user submit phone
func (h *UserHandler) CreateOTP(c echo.Context) error {
	temp := &struct {
		Phone string `json:"phone" form:"phone" query:"phone"`
	}{}
	if err := c.Bind(temp); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}
	if !validatePhone(temp.Phone) {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Phone is not valid"})
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	otp, err := h.UUsecase.CreateOTP(ctx, temp.Phone)
	if err != nil {
		fmt.Println("err 1:", err)
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":         otp.ID,
		"phone":      otp.Phone,
		"ticket":     otp.Ticket,
		"created_at": otp.CreatedAt,
		"expired_at": otp.ExpiredAt,
	})
}

// VerifyOTP ...
func (h *UserHandler) VerifyOTP(c echo.Context) (err error) {
	type tempStruct struct {
		OTP    string `json:"otp" form:"otp" query:"otp" validate:"required"`
		Phone  string `json:"phone" form:"phone" query:"phone" validate:"required"`
		Ticket string `json:"ticket" form:"ticket" query:"ticket" validate:"required"`
	}
	temp := new(tempStruct)
	if err = c.Bind(temp); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	if ok, err := isOTPRequestValid(temp); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if !validatePhone(temp.Phone) {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Phone is not valid"})
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	// find phone if exists, we do not need verify otp
	u, _ := h.UUsecase.GetByPhone(ctx, temp.Phone)

	if u != nil {
		// found phone number
		return c.JSON(http.StatusBadRequest, "Phone has already existed")
	}

	newUser := &models.User{
		Phone: temp.Phone,
	}
	if err = h.UUsecase.Store(ctx, newUser); err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	if err = h.UUsecase.VerifyOTP(ctx, temp.Phone, temp.OTP, temp.Ticket); err != nil {
		return c.JSON(http.StatusNotFound, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "otp verification is success"})
}

// GetByPhone ...
func (h *UserHandler) GetByPhone(c echo.Context) error {
	phone := c.QueryParam("phone")

	if !validatePhone(phone) {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Phone is not valid"})
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	art, err := h.UUsecase.GetByPhone(ctx, phone)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, art)
}

// Store will store the article by given request body
func (h *UserHandler) Store(c echo.Context) error {
	var user models.User
	err := c.Bind(&user)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	if ok, err := isRequestValid(&user); !ok {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = h.UUsecase.Store(ctx, &user)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	newUser, _ := h.UUsecase.GetByPhone(ctx, user.Phone)

	return c.JSON(http.StatusCreated, newUser)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	logrus.Error(err)
	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrConflict:
		return http.StatusConflict
	case models.ErrBadParamInput:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
