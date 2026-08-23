// Package util holds small, reusable primitives shared across common/.
package util

import (
	"io"
	"sync"
	"time"
)

// NewTimeoutReader wraps an io.Reader and calls kill when the stream
// goes quiet for too long. Every successful Read resets the idle
// timer; when the timer fires, or when the returned reader is Closed,
// kill is invoked. A non-positive timeout disables the watcher.
//
// The second return value is a callback that stashes an error on the
// reader -- only the source (caller of NewTimeoutReader) receives it,
// so external consumers can't overwrite the stashed error.
func NewTimeoutReader(r io.ReadCloser, kill func(), timeout time.Duration) (io.ReadCloser, func(error)) {
	if timeout <= 0 {
		return r, func(error) {}
	}

	tr := &TimeoutReader{
		ReadCloser: r,
		timeout:    timeout,
		reset:      make(chan any, 1),
		done:       make(chan struct{}),
	}
	go tr.watch(kill)
	return tr, tr.SetErr
}

// TimeoutReader wraps a ReadCloser and kills the underlying producer
// if the stream goes quiet for too long. Every successful Read resets
// the idle timer; when the timer fires the kill func is called, and
// the same happens on Close so the producer doesn't linger.
type TimeoutReader struct {
	io.ReadCloser
	timeout time.Duration
	reset   chan any
	done    chan struct{}
	once    sync.Once
	err     error
	errOnce sync.Once
}

// Err returns the stashed error from a timeout or non-zero exit.
// Returns nil if the stream completed normally.
func (t *TimeoutReader) Err() error {
	return t.err
}

// SetErr stashes an error (e.g. from a non-zero exit) so callers can
// retrieve it via Err after the stream ends.
func (t *TimeoutReader) SetErr(e error) {
	t.errOnce.Do(func() { t.err = e })
}

func (t *TimeoutReader) Read(p []byte) (int, error) {
	select {
	case t.reset <- nil:
	default:
	}
	return t.ReadCloser.Read(p)
}

func (t *TimeoutReader) Close() error {
	t.once.Do(func() {
		close(t.done)
	})
	return t.ReadCloser.Close()
}

func (t *TimeoutReader) watch(kill func()) {
	timer := time.NewTimer(t.timeout)
	defer timer.Stop()
	for {
		select {
		case <-t.reset:
			if !timer.Stop() { // if the timer already stopped, drain the channel.
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(t.timeout)
		case <-timer.C:
			t.SetErr(&TimeoutError{timeout: t.timeout})
			kill()
			return
		case <-t.done:
			return
		}
	}
}

// TimeoutError is returned by TimeoutReader when the idle timer fires.
type TimeoutError struct {
	timeout time.Duration
}

func (e *TimeoutError) Error() string {
	return "stream timed out after " + e.timeout.String()
}
