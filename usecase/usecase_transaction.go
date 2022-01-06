package usecase

import (
	"finance/adaptor"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type listTransactionResponse struct {
	Transactions []adaptor.Transaction `json:"transactions"`
	Message      string                `json:"message"`
}

func ListTransaction(c echo.Context) error {
	walletID, _ := strconv.Atoi(c.QueryParam("wallet_id"))

	user := c.Get("User").(*adaptor.User)
	owner, err := IsUserOwnWalletID(user.UserID, walletID)
	if err != nil {
		fmt.Printf("%v\n", err)
		return c.String(http.StatusBadRequest, `{"message": "Error Server"}`)
	}
	if !owner {
		return c.String(http.StatusUnauthorized, `{"message": "Error Permission"}`)
	}

	txs, err := transactionAdaptor.FilterTransaction(walletID, time.Now(), time.Now())
	if err != nil {
		fmt.Printf("%v\n", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Try again later"})
	}

	res := listTransactionResponse{
		Transactions: txs,
		Message:      "Success",
	}

	return c.JSON(http.StatusOK, res)
}

type addTransactionRequest struct {
	WalletID   int    `json:"wallet_id" validate:"required"`
	CategoryID int    `json:"category_id" validate:"required,gt=0"`
	Amount     int    `json:"amount" validate:"required,gt=0"`
	Note       string `json:"note"`
}

type addTransactionResponse struct {
	Transaction adaptor.Transaction `json:"transaction"`
	Message     string              `json:"message"`
}

func AddTransaction(c echo.Context) error {
	requestBody := new(addTransactionRequest)
	err := c.Bind(&requestBody)
	if err != nil {
		fmt.Printf("%v\n", err)
		return c.String(http.StatusBadRequest, `{"message": "Error input"}`)
	}

	err = c.Validate(requestBody)
	if err != nil {
		fmt.Printf("%v\n", err)
		return c.String(http.StatusBadRequest, `{"message": "Error input"}`)
	}

	user := c.Get("User").(*adaptor.User)
	var owner bool
	owner, err = IsUserOwnWalletID(user.UserID, requestBody.WalletID)
	if err != nil {
		return c.String(http.StatusBadRequest, `{"message": "Error Server"}`)
	}
	if !owner {
		return c.String(http.StatusUnauthorized, `{"message": "Error Permission"}`)
	}

	transactionID, err := transactionAdaptor.CreateTransaction(requestBody.WalletID, requestBody.CategoryID,
		requestBody.Amount, requestBody.Note)
	if err != nil {
		fmt.Printf("[Error - Transaction Usecase] Add Transaction: %v", err)
		return c.String(http.StatusInternalServerError, `{"message": "Error Server"}`)
	}

	tx, err := transactionAdaptor.GetTransactionByID(int(transactionID))
	if err != nil {
		fmt.Printf("[Error - Transaction Usecase] Add Transaction: %v", err)
		return c.String(http.StatusInternalServerError, `{"message": "Error Server"}`)
	}
	res := addTransactionResponse{
		Transaction: *tx,
		Message:     "Success",
	}
	return c.JSON(http.StatusOK, res)
}

type updateTransactionRequest struct {
	WalletID   int    `json:"wallet_id" validate:"required"`
	CategoryID int    `json:"category_id" validate:"required"`
	Amount     int    `json:"amount" validate:"gt=0"`
	Note       string `json:"note"`
}

func UpdateTransaction(c echo.Context) error {
	transactionID, err := strconv.Atoi(c.Param("transaction_id"))
	if err != nil {
		return c.String(http.StatusBadRequest, `{"message": "Error Input"}`)
	}

	requestBody := new(updateTransactionRequest)
	err = c.Bind(&requestBody)
	if err != nil {
		return c.String(http.StatusBadRequest, `{"message": "Error Input"}`)
	}

	// Validate user owns transactionID
	var own bool
	user := c.Get("User").(*adaptor.User)
	own, err = IsUserOwnTransactionID(user.UserID, transactionID)
	if err != nil {
		return c.String(http.StatusBadRequest, `{"message": "Error Server"}`)
	}
	if !own {
		return c.String(http.StatusUnauthorized, `{"message": "Error Permission"}`)
	}

	// Validate user owns wallet
	own, err = IsUserOwnWalletID(user.UserID, requestBody.WalletID)
	if err != nil {
		return c.String(http.StatusBadRequest, `{"message": "Error Server"}`)
	}
	if !own {
		return c.String(http.StatusUnauthorized, `{"message": "Error Permission"}`)
	}

	err = transactionAdaptor.UpdateTransaction(transactionID, requestBody.WalletID, requestBody.CategoryID,
		requestBody.Amount, requestBody.Note)
	if err != nil {
		return c.String(http.StatusUnauthorized, `{"message": "Error Server"}`)
	}

	return c.String(http.StatusOK, `{"message": "Success"}`)
}
