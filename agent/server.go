package agent

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type server struct {
	srv *http.Server
}

func newServer(listen string, handler http.Handler) *server {
	return &server{
		srv: &http.Server{
			Addr:    listen,
			Handler: handler,
		},
	}
}

func (s *server) listenAndServe() error {
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("http server: %w", err)
	}
	return nil
}

func (s *server) shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.srv.Shutdown(shutdownCtx)
}
