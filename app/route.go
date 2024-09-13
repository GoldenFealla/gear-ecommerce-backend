package app

import (
	gc "github.com/goldenfealla/gear-manager/module/gear/controller"
	uc "github.com/goldenfealla/gear-manager/module/user/controller"
	"github.com/labstack/echo/v4"
)

func Route(e *echo.Echo) {
	// Gear
	gear := e.Group("gear")
	gear.GET("/test", gc.Test)

	// User
	user := e.Group("user")
	user.GET("/test", uc.Test)
}
