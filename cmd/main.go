package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"

	"github.com/goldenfealla/gear-manager/config"
	"github.com/goldenfealla/gear-manager/internal/repository/postgres"
	"github.com/goldenfealla/gear-manager/internal/rest"
	"github.com/goldenfealla/gear-manager/usecase"
)

func main() {
	// Loading config
	c := config.Load()

	// Connect to database PostgreSQL
	conn, err := pgx.Connect(context.Background(), c.Postgres)

	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close(context.Background())

	// Connect to database Redis
	ropts, err := redis.ParseURL(c.Redis)
	if err != nil {
		log.Fatalln(err)
	}
	rdb := redis.NewClient(ropts)

	// init Session Store
	store, err := redisstore.NewRedisStore(context.Background(), rdb)
	if err != nil {
		log.Fatal("failed to create redis store: ", err)
	}

	e := echo.New()

	// Middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     c.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Range", "X-Content-Range"},
		Skipper:          middleware.DefaultSkipper,
		MaxAge:           3600,
	}))
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Request timeout",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			log.Println(c.Path())
		},
		Timeout: 10 * time.Second,
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		CustomTimeFormat: "15:04:05 02/01/2006",
		Format:           "time: ${time_custom}\tmethod=${method}\turi=${uri}\tstatus=${status}\tlatency=${latency}\terror=${error}\n",
	}))
	e.Use(session.Middleware(store))

	// set up validator
	v := validator.New(validator.WithRequiredStructEnabled())

	// Build Repository
	gr := postgres.NewGearRepository(conn)
	ur := postgres.NewUserRepository(conn)

	// Build Usecase
	gu := usecase.NewGearUsecase(gr)
	uu := usecase.NewUserUsecase(ur)

	// Build Handler
	rest.NewUserHandler(e, uu, v)
	rest.NewGearHandler(e, gu, v)

	err = e.Start(fmt.Sprintf("%v:%v", c.Host, c.Port))
	if err != nil {
		log.Fatalln(err)
	}
}
