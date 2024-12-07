package server

import (
	"kaf-interface/internal/orders/config"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + cfg.HTTP.Port,
			Handler: handler,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}
