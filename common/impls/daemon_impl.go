package impls

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

type DaemonImpl struct {
	mainFn func(ctx context.Context) error

	mu     sync.Mutex
	cancel context.CancelFunc
	ctx    context.Context

	errValue atomic.Value
}

func NewDaemonImpl(main func(ctx context.Context) error) *DaemonImpl {
	return &DaemonImpl{mainFn: main}
}

func (d *DaemonImpl) Start(onExit ...func()) error {
	if !d.mu.TryLock() {
		return errors.New("daemon: already running")
	}

	d.errValue.Store(nil)

	ctx, cancel := context.WithCancel(context.Background())
	d.cancel = cancel
	d.ctx = ctx

	callbacks := onExit

	go func() {
		defer d.mu.Unlock()
		defer func() {
			for _, fn := range callbacks {
				fn()
			}
		}()

		if err := d.mainFn(ctx); err != nil {
			d.errValue.Store(err)
		}
	}()

	return nil
}

func (d *DaemonImpl) Stop() error {
	if d.cancel == nil {
		return errors.New("daemon: not running")
	}

	d.cancel()
	return nil
}

func (d *DaemonImpl) Running() bool {
	if d.ctx == nil {
		return false
	}

	select {
	case <-d.ctx.Done():
		return false
	default:
		return true
	}
}

func (d *DaemonImpl) Error() error {
	if v := d.errValue.Load(); v != nil {
		return v.(error)
	}
	return nil
}
