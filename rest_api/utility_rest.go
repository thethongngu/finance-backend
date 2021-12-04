package rest_api

import (
	"github.com/labstack/echo/v4"
)

func readCookie(c echo.Context, name string) (string, error) {
	cookie, err := c.Cookie(name)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}
