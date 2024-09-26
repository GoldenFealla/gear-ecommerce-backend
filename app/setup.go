package app

import (
	"github.com/goldenfealla/gear-manager/internal/middleware"
	"github.com/goldenfealla/gear-manager/internal/rest"
)

func (s *Server) Setup() {
	s.e.Use(middleware.CORS)

	// Build Handler
	rest.NewUserHandler(s.e, nil)
	rest.NewGearHandler(s.e, nil)
}
