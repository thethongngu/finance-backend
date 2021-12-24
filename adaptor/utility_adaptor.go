package adaptor

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

// Reuse connection or create new connection depend on the logic
// Current version create new connection for each usecase (not each request)
func GetMySQLConnection() *sql.DB {
	var err error
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	cfg := mysql.Config{
		User:   username,
		Passwd: password,
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "finance",
	}
	mysql, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		fmt.Print("Error database...")
		panic(err)
	}

	mysql.SetConnMaxLifetime(time.Minute * 3)
	mysql.SetMaxOpenConns(10)
	mysql.SetMaxIdleConns(10)

	return mysql
}
