package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

	_otpRepo "github.com/lyxuansang91/credify-test/otp/repository"
	_userHttpDelivery "github.com/lyxuansang91/credify-test/user/delivery/http"
	_userRepo "github.com/lyxuansang91/credify-test/user/repository"
	_userUcase "github.com/lyxuansang91/credify-test/user/usecase"
)

// Init ...
func Init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

func urlSkipper(c echo.Context) bool {
	if strings.HasPrefix(c.Path(), "/metrics") {
		return true
	}
	return false
}

func main() {
	Init()
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sqlx.Connect(`mysql`, dsn)

	// dbConn, err := sql.Open(`mysql`, dsn)
	if err != nil && viper.GetBool("debug") {
		fmt.Println(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	p := prometheus.NewPrometheus("echo", urlSkipper)
	p.Use(e)
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	otpRepo := _otpRepo.NewMysqlOTPRepository(dbConn)
	// authorRepo := _authorRepo.NewMysqlAuthorRepository(dbConn)
	userRepo := _userRepo.NewMysqlUserRepository(dbConn)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	uu := _userUcase.NewUserUsecase(userRepo, otpRepo, timeoutContext)
	_userHttpDelivery.NewUserHandler(e, uu)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
