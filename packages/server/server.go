package server

import (
	"context"
	"net/http"
	"time"

	"github.com/Petr0752/L0/packages/config"
)

const (
	maxHeaderBytes = 1 << 20 // 1 MB
	timeout        = 10 * time.Second
)

type (
	Server struct {
		httpServer *http.Server
	}
)

func New(handler http.Handler) *Server {
	cfg := config.Get()
	s := &Server{}
	s.httpServer = &http.Server{
		Addr:           cfg.ServerAddress,
		Handler:        handler,
		MaxHeaderBytes: maxHeaderBytes,
		ReadTimeout:    timeout,
		WriteTimeout:   timeout,
	}

	return s
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
