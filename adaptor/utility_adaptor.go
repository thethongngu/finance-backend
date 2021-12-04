package adaptor

import (
	"database/sql"
	"time"
)

// Reuse connection or create new connection depend on the logic
// Current version create new connection for each usecase (not each request)
func GetMySQLConnection() *sql.DB {
	var err error
	mysql, err := sql.Open("mysql", "root:localroot@/finance") // TODO: get from env var
	if err != nil {
		panic(err)
	}

	mysql.SetConnMaxLifetime(time.Minute * 3)
	mysql.SetMaxOpenConns(10)
	mysql.SetMaxIdleConns(10)

	return mysql
}
