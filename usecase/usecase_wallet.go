package usecase

import (
	"finance/adaptor"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ListWallet(c echo.Context) error {
	user := c.Get("User").(*adaptor.User)
	wallets, err := walletAdaptor.GetWalletByUserID(user.UserID)
	if err != nil {
		return c.String(http.StatusUnauthorized, `{"message": "Error Server"}`)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"wallets": wallets,
		"message": "Success",
	})
}

func GetCurrency(c echo.Context) error {
	currencies, err := walletAdaptor.GetAllCurrency()
	if err != nil {
		return c.String(http.StatusUnauthorized, `{"message": "Error Server"}`)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"currencies": currencies,
		"message":    "Success",
	})
}
