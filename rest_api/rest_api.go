package rest_api

import (
	"finance/adaptor"
	"finance/usecase"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

var userUsecase *usecase.UserUsecase

func StartRESTAPIServer() {

	userAdaptor := adaptor.NewUserMySQLAdaptor()
	userUsecase = usecase.NewUserUsecase(userAdaptor)

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	// Add middware CORS, Logging

	e.Static("/static", "static")

	e.POST("/login", HandleLogin)
	e.Group("/member", ValidateUser)

	// e.POST("/transaction", addTransaction)
	// e.PUT("/transaction/:transaction_id", updateTransaction)
	// e.DELETE("/transaction/:transaction_id", deleteTransaction)

	e.Logger.Fatal(e.Start(":2808"))
}
