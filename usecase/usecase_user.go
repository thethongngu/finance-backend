package usecase

import (
	"finance/adaptor"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type userResponse struct {
	Username  string `json:"username"`
	UserID    int    `json:"user_id"`
	AvatarURL string `json:"avatar_url"`
}

type loginRequest struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
}

type loginResponse struct {
	User    userResponse `json:"user"`
	Message string       `json:"message"`
}

func HandleLogin(c echo.Context) error {
	var req loginRequest
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, echo.Map{"message": "Login data not enough"})
	}

	err = c.Validate(req)
	if err != nil {
		return c.String(http.StatusBadRequest, `{"message": "Error input"}`)
	}

	var user *adaptor.User
	user, err = userAdaptor.GetUserByCredentails(req.Username, req.HashedPassword)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Error wrong username or password"})
	}

	var session *adaptor.Session
	session, err = userAdaptor.CreateNewUserSession(user)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Error server"})
	}

	cookie := new(http.Cookie)
	cookie.Name = "session_id"
	cookie.Value = session.SessionID
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true
	cookie.Path = "/"
	c.SetCookie(cookie)

	userData := userResponse{
		Username: user.Username, UserID: user.UserID, AvatarURL: user.AvatarURL,
	}
	resData := loginResponse{User: userData, Message: "Login success"}
	return c.JSON(http.StatusOK, resData)
}

func HandleRemember(c echo.Context) error {
	sessionID, err := ReadCookie(c, "session_id")
	if err != nil {
		return c.String(http.StatusUnauthorized, `{"message": "Login failed"}`)
	}

	var user *adaptor.User
	user, err = userAdaptor.GetUserBySessionID(sessionID)
	if err != nil {
		return c.String(http.StatusUnauthorized, `{"message": "Login failed"}`)
	} else {
		c.Set("User", user)
	}

	userData := userResponse{
		Username: user.Username, UserID: user.UserID, AvatarURL: user.AvatarURL,
	}
	resData := loginResponse{User: userData, Message: "Login success"}
	return c.JSON(http.StatusOK, resData)
}

func ValidateUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionID, err := ReadCookie(c, "session_id")
		if err != nil {
			return c.String(http.StatusUnauthorized, `{"message": "Login failed"}`)
		}

		var user *adaptor.User
		user, err = userAdaptor.GetUserBySessionID(sessionID)
		if err != nil {
			return c.String(http.StatusUnauthorized, `{"message": "Login failed"}`)
		} else {
			c.Set("User", user)
			return next(c)
		}
	}
}
