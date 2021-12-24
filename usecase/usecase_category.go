package usecase

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetCategory(c echo.Context) error {
	category, err := categoryAdaptor.GetAllCategory()
	if err != nil {
		c.JSON(http.StatusBadRequest, echo.Map{"message": "Error Server"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"category": category,
		"message":  "Success",
	})
}
