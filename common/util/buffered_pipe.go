package util

import (
	"context"
	"io"
	"sync"
)

// BufferedPipe is a simple in-memory pipe with an optional capacity cap.
// When maxBuf is zero the pipe grows without bound; otherwise Write
// blocks once the buffer reaches maxBuf until Read drains it.
//
// Cancellation of the supplied context signals both readers and writers
// to unblock. A zero-value BufferedPipe is ready to use after calling Init.
type BufferedPipe struct {
	maxBuf int
	buf    []byte
	mu     sync.Mutex
	cond   *sync.Cond
	ctx    context.Context
	cancel context.CancelFunc
}

// BufferedPipeOptions configures a BufferedPipe.
type BufferedPipeOptions struct {
	// Ctx is the context that, when cancelled, closes the pipe. If nil,
	// a background context with a new cancel function is created.
	Ctx context.Context

	// Partition between mandatory and optional fields.

	// MaxBuf caps the buffer size. Zero or nil means unlimited.
	MaxBuf *int
}

// NewBufferedPipe returns a ReadCloser and a WriteCloser for an in-memory
// pipe configured with the given options.
//
// When MaxBuf is explicitly set to 0 (non-nil pointer to zero), a standard
// io.Pipe is used. Otherwise a capacity-capped BufferedPipe is returned.
func NewBufferedPipe(opts ...BufferedPipeOptions) (io.ReadCloser, io.WriteCloser) {
	var opt BufferedPipeOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.MaxBuf != nil && *opt.MaxBuf == 0 {
		r, w := io.Pipe()
		return r, w
	}

	p := &BufferedPipe{}
	srcCtx := opt.Ctx
	if srcCtx == nil {
		srcCtx = context.Background()
	}
	p.ctx, p.cancel = context.WithCancel(srcCtx)
	if opt.MaxBuf != nil {
		p.maxBuf = *opt.MaxBuf
	}
	p.cond = sync.NewCond(&p.mu)
	return p, p
}

// Write appends data to the pipe. If the pipe is full (maxBuf > 0) it
// blocks until space is freed by Read or the context is cancelled.
// Returns io.ErrClosedPipe if the context was cancelled before the
// write completed.
func (p *BufferedPipe) Write(data []byte) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.buf = append(p.buf, data...)
	if p.maxBuf > 0 {
		for len(p.buf) > p.maxBuf {
			if p.ctx.Err() != nil {
				return 0, io.ErrClosedPipe
			}
			p.cond.Wait()
		}
	}
	p.cond.Signal()
	return len(data), nil
}

// Read copies data from the pipe into dst. Blocks if the pipe is empty
// until data arrives or the context is cancelled. Returns io.EOF when
// the context is cancelled and the buffer is drained.
func (p *BufferedPipe) Read(dst []byte) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for len(p.buf) == 0 {
		if p.ctx.Err() != nil {
			return 0, io.EOF
		}
		p.cond.Wait()
	}

	bytesRead := copy(dst, p.buf)
	p.buf = p.buf[bytesRead:]
	p.cond.Signal()
	return bytesRead, nil
}

// Close cancels the pipe's context and wakes any waiting readers or writers.
// After Close, Read returns io.EOF and Write returns io.ErrClosedPipe.
func (p *BufferedPipe) Close() error {
	if p.cancel != nil {
		p.cancel()
	}
	p.cond.Broadcast()
	return nil
}
