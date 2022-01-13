package api

import (
	"finance/usecase"
	"fmt"
	"net/http"
	"os"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.validator.Struct(i)
	return err
}

func StartRESTAPIServer() {

	usecase.InitUsecase()

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	ip := os.Getenv("LOCAL_IP")
	port := os.Getenv("LOCAL_PORT")
	localAddress := "https://" + ip + ":" + port

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://localhost:3000", localAddress, "https://finance.namdeo.one"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.GET("/api/ping", func(c echo.Context) error {
		fmt.Println("ok")
		return c.String(http.StatusOK, "pong")
	})

	e.POST("/api/login", usecase.HandleLogin)
	e.POST("/api/remember", usecase.HandleRemember)

	e.GET("/api/currency", usecase.GetCurrency)

	member := e.Group("/api/member", usecase.ValidateUser)
	member.Static("/img", "img")

	member.GET("/wallet", usecase.ListWallet)

	member.GET("/transaction", usecase.ListTransaction)
	member.POST("/transaction", usecase.AddTransaction)
	member.PUT("/transaction/:transaction_id", usecase.UpdateTransaction)

	member.GET("/category", usecase.GetCategory)

	member.GET("/stats", usecase.GetStats)

	e.Logger.Fatal(e.Start(":2808"))
}
