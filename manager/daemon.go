package manager

import (
	"context"

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
	<-ctx.Done()
	return nil
}
