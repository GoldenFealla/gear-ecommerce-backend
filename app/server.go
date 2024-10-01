/*
Package app for initialize server and config

Use

	app.New()

to create a server
*/
package app

import (
	"context"
	"fmt"

	"github.com/goldenfealla/gear-manager/config"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type Server struct {
	e    *echo.Echo
	c    *config.Config
	conn *pgx.Conn
}

func New() *Server {
	c := config.Load()

	return &Server{
		e: echo.New(),
		c: c,
	}
}

/*
This function is the entry point, it will return any error for the main function
*/
func (s *Server) Start() error {
	err := s.connectToPostgres()
	if err != nil {
		return err
	}

	// Add handler
	s.Setup()

	err = s.e.Start(fmt.Sprintf("%v:%v", s.c.Host, s.c.Port))
	if err != nil {
		return err
	}

	defer s.conn.Close(context.Background())

	return nil
}

func (s *Server) connectToPostgres() error {
	conn, err := pgx.Connect(context.Background(), s.c.Postgres)

	if err != nil {
		return err
	}

	s.conn = conn
	return nil
}
