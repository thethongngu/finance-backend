package main

import (
	"finance/api"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	api.StartRESTAPIServer()
}
