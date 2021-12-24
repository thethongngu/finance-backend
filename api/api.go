package api

import (
	"finance/usecase"
	"fmt"
	"net/http"

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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://172.18.109.205:3000"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	e.Static("/static", "static")

	e.GET("/ping", func(c echo.Context) error {
		fmt.Println("ok")
		return c.String(http.StatusOK, "pong")
	})

	e.POST("/login", usecase.HandleLogin)
	e.POST("/remember", usecase.HandleRemember)

	e.GET("/currency", usecase.GetCurrency)

	member := e.Group("/member", usecase.ValidateUser)

	member.GET("/wallet", usecase.ListWallet)

	member.POST("/transaction", usecase.AddTransaction)
	member.PUT("/transaction/:transaction_id", usecase.UpdateTransaction)

	member.GET("/category", usecase.GetCategory)

	member.GET("/stats", usecase.GetStats)

	e.Logger.Fatal(e.Start(":2808"))
}
