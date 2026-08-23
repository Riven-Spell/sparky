package util

import (
	"context"
	"io"
	"testing"
)

func TestBufferedPipe_BasicReadWrite(t *testing.T) {
	readEnd, writeEnd := NewBufferedPipe()
	data := []byte("hello world")

	if _, err := writeEnd.Write(data); err != nil {
		t.Fatalf("Write: %v", err)
	}

	buf := make([]byte, 1024)
	n, err := readEnd.Read(buf)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if string(buf[:n]) != "hello world" {
		t.Fatalf("got %q, want %q", buf[:n], "hello world")
	}
}

func TestBufferedPipe_UnlimitedBuffer(t *testing.T) {
	readEnd, writeEnd := NewBufferedPipe()

	for i := 0; i < 1000; i++ {
		if _, err := writeEnd.Write([]byte("x")); err != nil {
			t.Fatalf("Write %d: %v", i, err)
		}
	}

	buf := make([]byte, 1024)
	n, err := readEnd.Read(buf)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if n != 1000 {
		t.Fatalf("got %d bytes, want 1000", n)
	}
}

func TestBufferedPipe_CappedBuffer(t *testing.T) {
	max := 5
	readEnd, writeEnd := NewBufferedPipe(BufferedPipeOptions{
		MaxBuf: &max,
	})

	if _, err := writeEnd.Write([]byte("hello")); err != nil {
		t.Fatalf("Write: %v", err)
	}

	errCh := make(chan error, 1)
	go func() {
		_, err := writeEnd.Write([]byte("world"))
		errCh <- err
	}()

	// Read to free space so the goroutine can complete.
	buf := make([]byte, 1024)
	n, err := readEnd.Read(buf)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if string(buf[:n]) != "hello" {
		t.Fatalf("got %q, want %q", buf[:n], "hello")
	}

	// Now the goroutine should finish.
	if err := <-errCh; err != nil {
		t.Fatalf("second Write: %v", err)
	}

	// Read the second chunk.
	n, err = readEnd.Read(buf)
	if err != nil {
		t.Fatalf("second Read: %v", err)
	}
	if string(buf[:n]) != "world" {
		t.Fatalf("got %q, want %q", buf[:n], "world")
	}
}

func TestBufferedPipe_ReadContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	readEnd, _ := NewBufferedPipe(BufferedPipeOptions{
		Ctx: ctx,
	})

	errCh := make(chan error, 1)
	go func() {
		buf := make([]byte, 1024)
		_, err := readEnd.Read(buf)
		errCh <- err
	}()

	cancel()

	err := <-errCh
	if err != io.EOF {
		t.Fatalf("got %v, want io.EOF", err)
	}
}

func TestBufferedPipe_WriteContextCancellation(t *testing.T) {
	max := 2
	ctx, cancel := context.WithCancel(context.Background())
	_, writeEnd := NewBufferedPipe(BufferedPipeOptions{
		Ctx:    ctx,
		MaxBuf: &max,
	})

	if _, err := writeEnd.Write([]byte("ab")); err != nil {
		t.Fatalf("first Write: %v", err)
	}

	errCh := make(chan error, 1)
	go func() {
		_, err := writeEnd.Write([]byte("cd"))
		errCh <- err
	}()

	// Give the goroutine time to block, then cancel.
	cancel()

	err := <-errCh
	if err != io.ErrClosedPipe {
		t.Fatalf("got %v, want io.ErrClosedPipe", err)
	}
}

func TestBufferedPipe_Close(t *testing.T) {
	readEnd, writeEnd := NewBufferedPipe()

	if _, err := writeEnd.Write([]byte("data")); err != nil {
		t.Fatalf("Write: %v", err)
	}

	if err := writeEnd.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	// Read buffered data first.
	buf := make([]byte, 1024)
	n, err := readEnd.Read(buf)
	if err != nil {
		t.Fatalf("first Read after Close: %v", err)
	}
	if string(buf[:n]) != "data" {
		t.Fatalf("got %q, want %q", buf[:n], "data")
	}

	// Now should get EOF.
	_, err = readEnd.Read(buf)
	if err != io.EOF {
		t.Fatalf("second Read after Close: got %v, want io.EOF", err)
	}
}

func TestBufferedPipe_NilContext(t *testing.T) {
	readEnd, writeEnd := NewBufferedPipe()

	if _, err := writeEnd.Write([]byte("test")); err != nil {
		t.Fatalf("Write: %v", err)
	}

	buf := make([]byte, 1024)
	n, err := readEnd.Read(buf)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if string(buf[:n]) != "test" {
		t.Fatalf("got %q, want %q", buf[:n], "test")
	}
}

func TestBufferedPipe_ConcurrentReadWrite(t *testing.T) {
	max := 10
	readEnd, writeEnd := NewBufferedPipe(BufferedPipeOptions{
		MaxBuf: &max,
	})

	errCh := make(chan error, 1)
	go func() {
		for i := 0; i < 100; i++ {
			if _, err := writeEnd.Write([]byte("x")); err != nil {
				errCh <- err
				return
			}
		}
		errCh <- nil
	}()

	buf := make([]byte, 1024)
	total := 0
	for total < 100 {
		n, err := readEnd.Read(buf)
		if err != nil {
			t.Fatalf("Read: %v", err)
		}
		total += n
	}

	if err := <-errCh; err != nil {
		t.Fatalf("Write: %v", err)
	}
}

func TestBufferedPipe_ExplicitZeroMaxBuf(t *testing.T) {
	zero := 0
	readEnd, writeEnd := NewBufferedPipe(BufferedPipeOptions{
		MaxBuf: &zero,
	})

	// io.Pipe: write blocks until read.
	errCh := make(chan error, 1)
	go func() {
		_, err := writeEnd.Write([]byte("hello"))
		errCh <- err
	}()

	buf := make([]byte, 1024)
	n, err := readEnd.Read(buf)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if string(buf[:n]) != "hello" {
		t.Fatalf("got %q, want %q", buf[:n], "hello")
	}

	if err := <-errCh; err != nil {
		t.Fatalf("Write: %v", err)
	}
}
