package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

type Transaction struct {
	Transaction_id int       `json:"transaction_id"`
	Wallet_id      int       `json:"wallet_id"`
	Category_id    int       `json:"category_id"`
	Amount         int       `json:"amount"`
	Created_at     time.Time `json:"created_at"`
	Note           string    `json:"note"`
}

func initialDBConnection() *sql.DB {
	db, err := sql.Open("mysql", "root:localroot@/finance")
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}

func handleHomePage(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func addTransaction(c echo.Context) error {
	body := echo.Map{}
	err := c.Bind(&body)
	if err != nil {
		fmt.Printf("%v", err)
	}

	_, err = db.Exec(`
		INSERT INTO transaction (wallet_id, category_id, amount, note) VALUES (?, ?, ?, ?)`,
		body["wallet_id"], body["category_id"], body["amount"], body["note"],
	)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return c.JSON(http.StatusOK, `{"message": "Added"}`)
}

func updateTransaction(c echo.Context) error {
	body := echo.Map{}
	err := c.Bind(&body)
	if err != nil {
		fmt.Printf("%v", err)
	}

	transaction_id, _ := strconv.Atoi(c.Param("transaction_id"))
	_, err = db.Exec(`
		UPDATE transaction SET wallet_id = ?, category_id = ?, amount = ?, note = ? 
		WHERE transaction_id = ?`,
		body["wallet_id"], body["category_id"], body["amount"], body["note"], transaction_id,
	)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return c.JSON(http.StatusOK, `{"message": "Updated"}`)
}

func deleteTransaction(c echo.Context) error {
	body := echo.Map{}
	err := c.Bind(&body)
	if err != nil {
		fmt.Printf("%v", err)
	}

	transaction_id, _ := strconv.Atoi(c.Param("transaction_id"))
	_, err = db.Exec(`
		DELETE FROM transaction WHERE transaction_id = ?`, transaction_id,
	)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return c.JSON(http.StatusOK, `{"message": "Deleted"}`)
}

var db *sql.DB

func main() {

	db = initialDBConnection()

	e := echo.New()
	e.GET("/", handleHomePage)
	e.POST("/transaction", addTransaction)
	e.PUT("/transaction/:transaction_id", updateTransaction)
	e.DELETE("/transaction/:transaction_id", deleteTransaction)
	e.Logger.Fatal(e.Start(":2808"))
}
