package middleware

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/goldenfealla/gear-manager/domain"
	"github.com/goldenfealla/gear-manager/internal/session"
	"github.com/labstack/echo/v4"
)

type AuthenticatedConfig struct {
	Excludes []string
}

func AuthenticatedWithConfig(co *AuthenticatedConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Try to get user info first
			userInfo, err := session.IsAuth(c)

			if userInfo != nil {
				c.Set("user", userInfo)
				return next(c)
			}

			if slices.Contains(co.Excludes, c.Path()) {
				return next(c)
			}

			return c.JSON(http.StatusUnauthorized, &domain.Response{
				Message: fmt.Sprintf("You are not logged in. Detail: %v", err.Error()),
			})
		}
	}
}
