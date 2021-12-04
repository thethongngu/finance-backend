package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func ValidateOwningTransaction(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		body := echo.Map{}
		err := c.Bind(&body)
		if err != nil {
			return c.String(http.StatusBadRequest, `{"message": "Error Input"}`)
		}

		var ok bool
		var wallet Wallet
		if wallet, ok = c.Get("Wallet").(Wallet); !ok {
			err = db.
				QueryRow(`SELECT * FROM transaction WHERE transaction_id = ?`, body["transaction_id"].(int)).
				Scan(&wallet.WalletID, &wallet.CurrencyID, &wallet.UserID)
			if err != nil {
				return c.String(http.StatusBadRequest, `{"message": "Error Input"}`)
			}
		}

		if wallet.WalletID != body["wallet_id"].(int) {
			return c.String(http.StatusUnauthorized, `{"message": "Error Permission"}`)
		}

		next(c)
	}
}

func addTransaction(c echo.Context) error {
	transaction := new(Transaction)
	err := c.Bind(&transaction)

	// Validate data presents or not
	if err != nil {
		fmt.Printf("%v\n", err)
		return c.String(http.StatusBadRequest, `{"message": "Error Input"}`)
	}

	// Validate fields (type, range)
	err = c.Validate(transaction)
	if err != nil {
		fmt.Printf("%v\n", err)
		return c.String(http.StatusBadRequest, `{"message": "Error Input"}`)
	}

	// Validate wallet ownership to add transaction
	var walletID int
	err = db.QueryRow(`SELECT wallet_id FROM wallet WHERE wallet_id = ?`, transaction.WalletID).Scan(&walletID)
	if err != nil {
		return c.String(http.StatusUnauthorized, `{"message": "Error Permission"}`)
	}
	user := c.Get("User").(User)
	if walletID != user.UserID {
		return c.String(http.StatusUnauthorized, `{"message": "Error Permission"}`)
	}

	// Add new transaction
	_, err = db.Exec(`
		INSERT INTO transaction (wallet_id, category_id, amount, note) VALUES (?, ?, ?, ?)`,
		transaction.WalletID, transaction.CategoryID, transaction.Amount, transaction.Note,
	)
	if err != nil {
		fmt.Printf("%v\n", err)
		return c.String(http.StatusInternalServerError, `{"message": "Error DB"}`)
	}

	return c.String(http.StatusOK, `{"message": "Added"}`)
}

func updateTransaction(c echo.Context) error {
	transaction := new(Transaction)
	err := c.Bind(&transaction)

	// Validate data fields present or not
	if err != nil {
		fmt.Printf("%v", err)
		return c.String(http.StatusBadRequest, `{"message": "Error input"}`)
	}

	// Validate fields (type, range)
	err = c.Validate(transaction)
	if err != nil {
		fmt.Printf("%v\n", err)
		return c.String(http.StatusBadRequest, `{"message": "Error input"}`)
	}

	// Validate user owns the transaction
	user := c.Get("User").(User)
	transactionIDURL, err := strconv.Atoi(c.Param("transaction_id"))
	if err != nil {
		return c.String(http.StatusNotFound, `{"message": "Error input"}`)
	}
	var currentWalletID int
	err = db.QueryRow(`SELECT wallet_id FROM transaction WHERE transaction_id = ?`, transactionIDURL).
	Scan(&currentWalletID)
	if err != nil {
		return c.String(http.StatusInternalServerError, `{"message": "Error DB"}`)
	}

	var currUserID int
	err = db.QueryRow(`SELECT user_id FROM wallet WHERE wallet_id = ?`, currentWalletID).Scan(&currUserID)
	if err != nil {
		return c.String(http.StatusInternalServerError, `{"message": "Error DB"}`)
	}
	
	if currUserID != user.UserID {
		return c.String(http.StatusUnauthorized, `{"message": "Error Permission"}`)
	}

	// Validate user own wallet_id in request body
	err = db.QueryRow(`SELECT user_id FROM wallet WHERE wallet_id = ?`, ).Scan(&transaction.WalletID)
	
	if currentWalletID != user. {
		return c.String(http.StatusUnauthorized, `{"message": "Error Permission"}`)
	}

	_, err = db.Exec(`
		UPDATE transaction SET wallet_id = ?, category_id = ?, amount = ?, note = ? WHERE transaction_id = ?`,
		transaction.WalletID, transaction.CategoryID, transaction.Amount, transaction.Note, transactionIDURL,
	)
	if err != nil {
		return c.String(http.StatusInternalServerError, `{"message": "Error DB"}`)
	}

	return c.JSON(http.StatusOK, `{"message": "Updated"}`)
}

// Request body: {}
// 1. Validate input data (range, value, type)
// 2. Validate permission (correct wallet_id, correct transaction_id of wallet_id)
// 3. Add new transaction
func deleteTransaction(c echo.Context) error {

	transaction_id, _ := strconv.Atoi(c.Param("transaction_id"))

	if err != nil {
		fmt.Printf("%v", err)
		return c.String(http.StatusInternalServerError, `{"message": "Error DB"}`)
	}

	isCorrect, errFunc := checkCorrectTransaction(c, transaction.TransactionID, transaction.WalletID)
	if !isCorrect {
		return errFunc()
	}

	isCorrect, errFunc = checkCorrectUser(c, transaction.WalletID, currUser.UserID)
	if !isCorrect {
		return errFunc()
	}

	_, err = db.Exec(`
		DELETE FROM transaction WHERE transaction_id = ?`, transaction_id,
	)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return c.JSON(http.StatusOK, `{"message": "Deleted"}`)
}
