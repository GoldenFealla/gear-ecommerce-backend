package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type GearUsecase interface{}

type GearHandler struct {
	uc GearUsecase
}

func NewGearHandler(e *echo.Echo, uc GearUsecase) {
	handler := &GearHandler{
		uc,
	}

	group := e.Group("gear")

	group.GET("/test", handler.Test)
}

func (h *GearHandler) Test(c echo.Context) error {
	return c.JSON(http.StatusOK, "Test gear Ok")
}
