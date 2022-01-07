package adaptor

import (
	"database/sql"
	"fmt"
	"time"
)

type Transaction struct {
	TransactionID int       `json:"transaction_id,omitempty"`
	WalletID      int       `json:"wallet_id" validate:"required"`
	CategoryID    int       `json:"category_id" validate:"required"`
	Amount        int       `json:"amount" validate:"required,gt=0"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	Note          string    `json:"note"`
}

type TransactionAdaptorInterface interface {
	GetTransactionByID(transactionID int) (*Transaction, error)
	FilterTransaction(walletID int, from time.Time, to time.Time) ([]Transaction, error)
	CreateTransaction(walletID int, categoryID int, amount int, note string) (int64, error)
	UpdateTransaction(transactionID int, walletID int, categoryID int, amount int, note string) error
	DeleteTransaction(transactionID int) error

	GetTotalAmount(walletID int, from time.Time, to time.Time) (int, error)
}

type TransactionMySQLAdaptor struct {
	conn *sql.DB
}

func NewTransactionMySQLAdaptor() TransactionMySQLAdaptor {
	return TransactionMySQLAdaptor{conn: GetMySQLConnection()}
}

func (a TransactionMySQLAdaptor) GetTransactionByID(transactionID int) (*Transaction, error) {
	var transaction Transaction
	err := a.conn.QueryRow(`SELECT * FROM Transaction WHERE transaction_id = ?`, transactionID).Scan(
		&transaction.TransactionID, &transaction.WalletID, &transaction.CategoryID,
		&transaction.Amount, &transaction.CreatedAt, &transaction.Note)
	if err != nil {
		err = fmt.Errorf("[Error] sql: %v", err)
		return nil, err
	}

	return &transaction, nil
}

func (a TransactionMySQLAdaptor) CreateTransaction(walletID int, categoryID int,
	amount int, note string) (int64, error) {

	result, err := a.conn.Exec(`
		INSERT INTO Transaction (wallet_id, category_id, amount, note) VALUES (?, ?, ?, ?)`,
		walletID, categoryID, amount, note,
	)
	if err != nil {
		err = fmt.Errorf("[Error] sql: %v", err)
		return -1, err
	}

	transactionID, _ := result.LastInsertId()
	return transactionID, nil
}

func (a TransactionMySQLAdaptor) UpdateTransaction(transactionID int, walletID int,
	categoryID int, amount int, note string) error {

	_, err := a.conn.Exec(`
		UPDATE Transaction SET wallet_id = ?, category_id = ?, amount = ?, note = ?) 
		WHERE transaction_id = ?`,
		transactionID, walletID, categoryID, amount, note,
	)
	if err != nil {
		err = fmt.Errorf("[Error] sql: %v", err)
		return err
	}

	return nil
}

func (a TransactionMySQLAdaptor) DeleteTransaction(transactionID int) error {
	return nil
}

func (a TransactionMySQLAdaptor) GetTotalAmount(walletID int, from time.Time, to time.Time) (int, error) {
	var total int
	err := a.conn.QueryRow(`
		SELECT coalesce(sum(amount), 0) FROM Transaction WHERE wallet_id = ? AND created_at BETWEEN ? AND ?`,
		walletID, from.Format("2006-01-02"), to.Format("2006-01-02"),
	).Scan(&total)
	if err != nil {
		return -1, err
	}

	return total, nil
}

func (a TransactionMySQLAdaptor) FilterTransaction(walletID int, from time.Time, to time.Time) ([]Transaction, error) {
	rows, err := a.conn.Query(`
		SELECT * FROM Transaction WHERE wallet_id = ? AND created_at BETWEEN ? AND ?`,
		walletID, from.Format("2006-01-02"), to.Format("2006-01-02"),
	)
	if err != nil {
		fmt.Printf("[Error - Transaction adaptor] Filter Transaction: %v", err)
		return nil, err
	}
	defer rows.Close()

	transactions := make([]Transaction, 0)
	for rows.Next() {
		var tx Transaction
		if err := rows.Scan(&tx.TransactionID, &tx.WalletID, &tx.CategoryID, &tx.Amount,
			&tx.CreatedAt, &tx.Note); err != nil {

			err = fmt.Errorf("[Error - Transaction adaptor] Filter Transaction:: %v", err)
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, err
}
