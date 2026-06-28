package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/LucasM4r/repomind/internal/api/router"
)

type Server struct {
	router  *router.Router
	address string
	port    string
}

func NewServer(router *router.Router, address, port string) *Server {
	return &Server{
		router:  router,
		address: address,
		port:    port,
	}
}

func (s *Server) Start() error {
	s.router.RegisterHandlers()

	addr := fmt.Sprintf("%s:%s", s.address, s.port)
	log.Printf("starting server on %s", addr)

	srv := &http.Server{
		Addr:         addr,
		Handler:      s.router.Handler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return srv.ListenAndServe()
}
