package app

import (
	"fmt"

	"github.com/goldenfealla/gear-manager/config"
	"github.com/labstack/echo/v4"
)

type Server struct {
	e *echo.Echo
	c *config.Config
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
	Route(s.e)

	err := s.e.Start(fmt.Sprintf("%v:%v", s.c.Host, s.c.Port))
	if err != nil {
		return err
	}

	return nil
}
