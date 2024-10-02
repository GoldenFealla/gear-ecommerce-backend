package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/goldenfealla/gear-manager/config"
	"github.com/goldenfealla/gear-manager/internal/repository/postgres"
	"github.com/goldenfealla/gear-manager/internal/rest"
	"github.com/goldenfealla/gear-manager/usecase"
)

func main() {
	// Loading config
	c := config.Load()

	// Connect to database
	conn, err := pgx.Connect(context.Background(), c.Postgres)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close(context.Background())

	e := echo.New()

	// Middleware
	e.Use(middleware.CORS())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Request timeout",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			log.Println(c.Path())
		},
		Timeout: 10 * time.Second,
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	// set up validator
	v := validator.New(validator.WithRequiredStructEnabled())

	// Build Repository
	gr := postgres.NewGearRepository(conn)

	// Build Usecase
	gu := usecase.NewGearUsecase(gr)

	// Build Handler
	rest.NewUserHandler(e, nil)
	rest.NewGearHandler(e, gu, v)

	err = e.Start(fmt.Sprintf("%v:%v", c.Host, c.Port))
	if err != nil {
		log.Fatalln(err)
	}
}
