package usecase

import "github.com/labstack/echo/v4"

func ReadCookie(c echo.Context, name string) (string, error) {
	cookie, err := c.Cookie(name)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func IsUserOwnWalletID(userID int, walletID int) (bool, error) {
	wallet, err := walletAdaptor.GetWalletByID(walletID)
	if err != nil {
		return false, err
	}
	return userID == wallet.UserID, nil
}

func IsUserOwnTransactionID(userID int, transactionID int) (bool, error) {
	wallet, err := walletAdaptor.GetWalletByTransactionID(transactionID)
	if err != nil {
		return false, err
	}

	return userID == wallet.UserID, nil
}
