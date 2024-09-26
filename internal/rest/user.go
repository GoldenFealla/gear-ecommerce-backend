package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserUsecase interface{}

type UserHandler struct {
	uc UserUsecase
}

func NewUserHandler(e *echo.Echo, uc UserUsecase) {
	handler := &UserHandler{
		uc,
	}

	group := e.Group("user")

	group.GET("/test", handler.Test)
}

func (h *UserHandler) Test(c echo.Context) error {
	return c.JSON(http.StatusOK, "Test user Ok")
}
