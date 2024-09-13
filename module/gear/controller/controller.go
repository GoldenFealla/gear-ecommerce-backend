package controller

import (
	"fmt"
	"net/http"

	"github.com/goldenfealla/gear-manager/module/gear/usecase"
	"github.com/labstack/echo/v4"
)

func Test(c echo.Context) error {
	r := usecase.Test()

	return c.JSON(http.StatusOK, fmt.Sprintf("Gear: %v", r))
}
