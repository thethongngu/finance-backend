module finance

go 1.17

// replace "github.com/thethongngu/finance/entity" => ./entities/entity
// replace "github.com/thethongngu/finance-backend/api" => ./apis/api

require (
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/go-sql-driver/mysql v1.6.0
	github.com/labstack/echo/v4 v4.6.1
)

require gopkg.in/go-playground/assert.v1 v1.2.1 // indirect

require (
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.9.0
	github.com/labstack/gommon v0.3.1 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	golang.org/x/crypto v0.0.0-20211108221036-ceb1ce70b4fa // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/sys v0.0.0-20211113001501-0c823b97ae02 // indirect
	golang.org/x/text v0.3.7 // indirect
)
