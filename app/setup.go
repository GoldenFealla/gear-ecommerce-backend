package app

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/goldenfealla/gear-manager/internal/repository/postgres"
	"github.com/goldenfealla/gear-manager/internal/rest"
	"github.com/goldenfealla/gear-manager/usecase"
)

func (s *Server) Setup() {
	s.e.Use(middleware.CORS())
	s.e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Request timeout",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			log.Println(c.Path())
		},
		Timeout: 10 * time.Second,
	}))
	s.e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	// set up validator
	v := validator.New(validator.WithRequiredStructEnabled())

	// Build Repository
	gr := postgres.NewGearRepository(s.conn)

	// Build Usecase
	gu := usecase.NewGearUsecase(gr)

	// Build Handler
	rest.NewUserHandler(s.e, nil)
	rest.NewGearHandler(s.e, gu, v)
}
