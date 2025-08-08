package server

import (
	"calculator-go/pkg/config"
	"context"
	"fmt"
	"net/http"
)

type Server struct {
	server *http.Server
}

func New(cfg config.ServerCfg, handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:           fmt.Sprintf("%v:%v", cfg.Host, cfg.Port),
			Handler:        handler,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
			MaxHeaderBytes: cfg.MaxHeaderByte,
		},
	}
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
