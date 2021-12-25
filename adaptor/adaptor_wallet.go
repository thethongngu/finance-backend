package adaptor

import (
	"database/sql"
	"fmt"
)

type Wallet struct {
	WalletID   int `json:"wallet_id"`
	CurrencyID int `json:"currency_id"`
	UserID     int `json:"user_id"`
}

type Currency struct {
	CurrencyID int    `json:"currency_id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
}

type WalletAdaptorInterface interface {
	GetWalletByTransactionID(transactionID int) (*Wallet, error)
	GetWalletByID(walletID int) (*Wallet, error)
	GetWalletByUserID(userID int) ([]Wallet, error)
	GetAllCurrency() ([]Currency, error)
}

type WalletMySQLAdaptor struct {
	conn *sql.DB
}

func NewWalletMySQLAdaptor() WalletMySQLAdaptor {
	return WalletMySQLAdaptor{conn: GetMySQLConnection()}
}

func (a WalletMySQLAdaptor) GetWalletByTransactionID(transactionID int) (*Wallet, error) {
	var walletID int
	err := a.conn.
		QueryRow(`SELECT wallet_id FROM transaction WHERE transaction_id = ?`, transactionID).
		Scan(walletID)
	if err != nil {
		err = fmt.Errorf("[Error] sql: %v", err)
		return nil, err
	}

	var wallet *Wallet
	wallet, err = a.GetWalletByID(walletID)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (a WalletMySQLAdaptor) GetWalletByID(walletID int) (*Wallet, error) {
	var wallet Wallet
	err := a.conn.
		QueryRow(`SELECT wallet_id, currency_id, user_id FROM wallet WHERE wallet_id = ?`, walletID).
		Scan(&wallet.WalletID, &wallet.CurrencyID, &wallet.UserID)
	if err != nil {
		err = fmt.Errorf("[Error] sql: %v", err)
		return nil, err
	}

	return &wallet, nil
}

func (a WalletMySQLAdaptor) GetWalletByUserID(userID int) ([]Wallet, error) {
	rows, err := a.conn.Query(`SELECT * FROM Wallet WHERE user_id = ?`, userID)
	if err != nil {
		err = fmt.Errorf("[Error] GetWalletByUserID: %v", err)
		return nil, err
	}
	defer rows.Close()

	var wallets []Wallet
	for rows.Next() {
		var wallet Wallet
		if err := rows.Scan(&wallet.WalletID, &wallet.CurrencyID, &wallet.UserID); err != nil {
			err = fmt.Errorf("[Error] GetWalletByUserID: %v", err)
			return nil, err
		}
		wallets = append(wallets, wallet)
	}

	return wallets, err
}

func (a WalletMySQLAdaptor) GetAllCurrency() ([]Currency, error) {
	rows, err := a.conn.Query(`SELECT * FROM Currency`)
	if err != nil {
		err = fmt.Errorf("[Error] GetAllCurrency: %v", err)
		return nil, err
	}
	defer rows.Close()

	var currencies []Currency
	for rows.Next() {
		var currency Currency
		if err := rows.Scan(&currency.CurrencyID, &currency.Code, &currency.Name); err != nil {
			err = fmt.Errorf("[Error] GetAllCurrency: %v", err)
			return nil, err
		}
		currencies = append(currencies, currency)
	}

	return currencies, err
}
