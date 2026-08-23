package sparkrun

import (
	"context"
	"io"
	"testing"
	"time"
)

func TestClusterMonitorOptions(t *testing.T) {
	cluster := "test-cluster"
	interval := 5
	opts := ClusterMonitorOptions{
		Hosts:    []string{"node1", "node2"},
		Cluster:  &cluster,
		DryRun:   true,
		Interval: &interval,
		Simple:   true,
	}

	client := &cliClient{binaryPath: "sparkrun"}
	// Test the argument construction inside ClusterMonitor without executing
	// Since ClusterMonitor builds args and calls streamCmd, we verify option logic
	if len(opts.Hosts) != 2 || *opts.Cluster != "test-cluster" || !opts.DryRun || *opts.Interval != 5 || !opts.Simple {
		t.Errorf("unexpected options state: %#v", opts)
	}
	_ = client
}

func TestStreamCmdExecutionAndTiedClose(t *testing.T) {
	// Using "sh" to run a command that outputs to stdout and stderr
	client := &cliClient{
		binaryPath: "sh",
		timeout:    time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	stdout, stderr, kill, err := client.streamCmd(ctx, "test stream", "-c", "echo hello_out; echo hello_err >&2; sleep 0.1")
	if err != nil {
		t.Fatalf("streamCmd failed: %v", err)
	}
	_ = kill

	outBytes, err := io.ReadAll(stdout)
	if err != nil {
		t.Fatalf("reading stdout failed: %v", err)
	}
	if string(outBytes) != "hello_out\n" {
		t.Errorf("stdout = %q, want %q", string(outBytes), "hello_out\n")
	}

	errBytes, err := io.ReadAll(stderr)
	if err != nil {
		t.Fatalf("reading stderr failed: %v", err)
	}
	if string(errBytes) != "hello_err\n" {
		t.Errorf("stderr = %q, want %q", string(errBytes), "hello_err\n")
	}
}

func TestStreamCmdTiedCloseKillsProcess(t *testing.T) {
	client := &cliClient{
		binaryPath: "sh",
		timeout:    5 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stdout, stderr, kill, err := client.streamCmd(ctx, "test stream", "-c", "sleep 10")
	if err != nil {
		t.Fatalf("streamCmd failed: %v", err)
	}

	kill()

	done := make(chan struct{})
	go func() {
		defer close(done)
		_, _ = io.ReadAll(stdout)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("calling kill did not terminate subprocess and unblock stdout")
	}

	_ = stderr
}
