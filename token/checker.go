package token

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Checker(token string) echo.MiddlewareFunc {
	auth := "OAuth " + token
	errorMessage := struct {
		Message string `json:"message"`
	}{
		Message: "invalid or missing token",
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token != auth {
				return c.JSON(http.StatusUnauthorized, errorMessage)
			}
			c.Response().Header().Set(echo.HeaderServer, "ferrite")
			return next(c)
		}
	}
}
