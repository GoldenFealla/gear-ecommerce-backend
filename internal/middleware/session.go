package middleware

import (
	"net/http"

	"github.com/goldenfealla/gear-manager/domain"
	"github.com/goldenfealla/gear-manager/internal/session"
	"github.com/labstack/echo/v4"
)

func IsAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userInfo, err := session.IsAuth(c)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
				Message: err.Error(),
			})
		}

		if userInfo != nil {
			c.Set("user", userInfo)
			return next(c)
		}

		return c.JSON(http.StatusUnauthorized, &domain.ResponseError{
			Message: "You are not logged in",
		})
	}
}

func IsNotAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userInfo, err := session.IsAuth(c)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, &domain.ResponseError{
				Message: err.Error(),
			})
		}

		if userInfo != nil {
			return c.JSON(http.StatusOK, &domain.ResponseError{
				Message: "You are already logged in",
				Data:    userInfo,
			})
		}

		return next(c)
	}
}
