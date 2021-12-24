package usecase

import (
	"finance/adaptor"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type statsRequest struct {
	WalletID int `json:"wallet_id"`
}

type statsResponse struct {
	AmountToday int     `json:"amount_today"`
	AmountMonth int     `json:"amount_month"`
	AvgMonth    float32 `json:"avg_month"`
	Message     string  `json:"message"`
}

func GetStats(c echo.Context) error {
	// var req statsRequest
	// err := c.Bind(&req)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, echo.Map{"message": "Error input"})
	// }

	walletID, _ := strconv.Atoi(c.QueryParam("wallet_id"))

	user := c.Get("User").(*adaptor.User)
	owner, err := IsUserOwnWalletID(user.UserID, walletID)
	if err != nil {
		return c.String(http.StatusBadRequest, `{"message": "Error Server"}`)
	}
	if !owner {
		return c.String(http.StatusUnauthorized, `{"message": "Error Permission"}`)
	}

	amountToday, err := transactionAdaptor.GetTotalAmount(walletID, time.Now(), time.Now())
	if err != nil {
		fmt.Print(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Try again later"})
	}

	t := time.Now()
	firstDay := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, 0).Add(time.Nanosecond * -1)
	amountMonth, err := transactionAdaptor.GetTotalAmount(walletID, firstDay, lastDay)
	if err != nil {
		fmt.Print(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Try again later"})
	}

	numDate := lastDay.Sub(firstDay).Hours() / 24

	res := statsResponse{
		AmountToday: amountToday,
		AmountMonth: amountMonth,
		AvgMonth:    float32(amountMonth) / float32(numDate),
		Message:     "Success",
	}
	return c.JSON(http.StatusOK, res)
}
