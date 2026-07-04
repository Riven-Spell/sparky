package agent

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Riven-Spell/sparky/common/impls"
)

type Daemon struct {
	impls.DaemonImpl
	config Config
}

func NewDaemon(cfg Config) *Daemon {
	d := &Daemon{config: cfg}
	d.DaemonImpl = *impls.NewDaemonImpl(d.mainLoop)
	return d
}

func (d *Daemon) mainLoop(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/health", handleGetHealth(d))

	srv := newServer(d.config.listenAddress(), mux)

	go func() {
		<-ctx.Done()
		_ = srv.shutdown(context.Background())
	}()

	go d.runHealthChecks(ctx)

	if err := srv.listenAndServe(); err != nil {
		return fmt.Errorf("agent server: %w", err)
	}

	return nil
}
