package rest_api

import (
	"fmt"
	"net/http"
	"time"

	"finance/entity"

	"github.com/labstack/echo/v4"
)

func ValidateUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionID, err := readCookie(c, "session_id")
		if err != nil {
			return c.String(http.StatusUnauthorized, `{"message": "Login failed"}`)
		}

		var user *entity.User
		user, err = userUsecase.LoginBySession(sessionID)
		if err != nil {
			return c.String(http.StatusUnauthorized, `{"message": "Login failed"}`)
		} else {
			c.Set("User", user)
			return next(c)
		}
	}
}

type loginRequest struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
}

func HandleLogin(c echo.Context) error {
	var req loginRequest
	err := c.Bind(&req)
	if err != nil {
		fmt.Printf("[Error]: %v", err)
		c.JSON(http.StatusBadRequest, echo.Map{"message": "Login data not enough"})
	}

	var session *entity.Session
	session, err = userUsecase.LoginByPassword(req.Username, req.HashedPassword)
	if err != nil {
		fmt.Printf("[Error]: %v", err)
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": err,
		})
	}

	cookie := new(http.Cookie)
	cookie.Name = "session_id"
	cookie.Value = session.SessionID
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true
	cookie.Path = "/"
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Login success",
	})
}
