package sparkrun

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/Riven-Spell/sparky/common/util"
)

// cliClient is the real [Client] implementation: it shells out to the
// `sparkrun` binary. Fields are unexported; configuration goes
// through [CliOption].
type cliClient struct {
	binaryPath string
	env        []string
	timeout    time.Duration
	maxBuf     *int
}

// CliOption configures a [cliClient] returned by [NewCliClient].
type CliOption func(*cliClient)

// CliClientOptions is a zero-value discoverability surface for the
// options accepted by [NewCliClient]. Each method returns
// a [CliOption] that can be passed to NewCliClient -- they are
// aliases for the package-level WithBinaryPath / WithEnv /
// WithStreamTimeout functions, so the two forms are interchangeable.
//
// Example:
//
//	client := sparkrun.NewCliClient(
//	    sparkrun.CliClientOptions{}.WithBinaryPath("/usr/local/bin/sparkrun"),
//	    sparkrun.CliClientOptions{}.WithStreamTimeout(2 * time.Minute),
//	)
//
// Callers that prefer the package-level functions can use them
// directly; the methods only exist to make option discovery easier.
type CliClientOptions struct{}

// WithBinaryPath overrides the executable path. Default: [DefaultBinary] (looked up via $PATH).
func (CliClientOptions) WithBinaryPath(path string) CliOption {
	return func(c *cliClient) { c.binaryPath = path }
}

// WithEnv appends extra environment variables to the subprocess.
// Values are passed alongside the parent's environment.
func (CliClientOptions) WithEnv(env []string) CliOption {
	return func(c *cliClient) { c.env = append(c.env, env...) }
}

// WithStreamBuffer sets the maximum buffer size for each stream (stdout
// and stderr) returned by [Client.Stream]. When nil, streams use an
// unbounded buffer (default). Capping the buffer prevents a slow reader
// from blocking the subprocess writer indefinitely.
func (CliClientOptions) WithStreamBuffer(size *int) CliOption {
	return func(c *cliClient) { c.maxBuf = size }
}

// WithStreamTimeout sets the maximum idle time between reads from a
// streamed subprocess. Each successful Read resets the timer; if the
// subprocess produces nothing for this long, streamCmd kills it. A
// zero value disables the timeout. Default: 1m.
func (CliClientOptions) WithStreamTimeout(d time.Duration) CliOption {
	return func(c *cliClient) { c.timeout = d }
}

// NewCliClient returns a [Client] backed by a real sparkrun
// subprocess. The returned interface lets callers depend on the
// contract rather than the implementation, so tests can substitute a
// fake later without changing call sites.
//
// Options are passed as [CliOption] values, pulled from methods on [CliClientOptions].
func NewCliClient(opts ...CliOption) Client {
	c := &cliClient{
		binaryPath: DefaultBinary,
		timeout:    time.Minute,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// runCmd executes the binary with args and returns stdout, stderr,
// and any error.
//
// The returned error is already classified: [ExecError] if the
// subprocess could not start, [ExitError] on a non-zero code, or
// [UsageError] for click's exit-2 invalid-flag case.
func (c *cliClient) runCmd(ctx context.Context, args ...string) (stdout, stderrBytes []byte, err error) {
	cmd := exec.CommandContext(ctx, c.binaryPath, args...)
	if len(c.env) > 0 {
		cmd.Env = append(os.Environ(), c.env...)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	waitErr := cmd.Run()
	stdout = stdoutBuf.Bytes()
	stderrBytes = stderrBuf.Bytes()
	if waitErr != nil {
		subcommand := subcommandPath(args)
		var exitErr *exec.ExitError
		if errors.As(waitErr, &exitErr) {
			// Click exits with code 2 for usage errors; anything else
			// is "sparkrun had something to say about it".
			if exitErr.ExitCode() == 2 {
				return stdout, stderrBytes, newUsageError(subcommand, c.binaryPath, args, exitErr.ExitCode(), string(stderrBytes), string(stdout), waitErr)
			}
			return stdout, stderrBytes, newExitError(subcommand, c.binaryPath, args, exitErr.ExitCode(), string(stderrBytes), string(stdout), waitErr)
		}

		return stdout, stderrBytes, newExecError(subcommand, c.binaryPath, args, string(stderrBytes), string(stdout), waitErr)
	}

	return stdout, stderrBytes, nil
}

// jsonCmd runs the command with --json appended and decodes stdout
// into out. Non-zero exits surface as [ExitError] / [UsageError]; the
// raw body is retained via StdError and StdOut so callers can inspect
// what sparkrun said. Empty/whitespace bodies decode as a no-op (out
// is left untouched).
//
// [ParseError] is reserved for the (rare) case where sparkrun exited
// cleanly but the body couldn't be decoded.
func (c *cliClient) jsonCmd(ctx context.Context, out any, args ...string) error {
	args = append(args, "--json")
	stdout, stderrBytes, err := c.runCmd(ctx, args...)
	if err != nil {
		return err
	}
	if out == nil {
		return nil
	}
	if len(bytes.TrimSpace(stdout)) == 0 {
		return nil
	}
	if err := json.Unmarshal(stdout, out); err != nil {
		return newParseError(subcommandPath(args), c.binaryPath, args, stdout, string(stderrBytes), err)
	}
	return nil
}

// plainCmd runs the command without --json. The wrapper captures the
// raw stdout into StdOut on non-zero exit so callers can log it, but
// otherwise discards it.
func (c *cliClient) plainCmd(ctx context.Context, args ...string) error {
	_, _, err := c.runCmd(ctx, args...)
	return err
}

// streamCmd runs args and returns two [io.Reader] streams (stdout and
// stderr) plus a [kill] function for hard-exiting the subprocess.
//
// The subprocess is terminated when the command exits, when the provided
// context is cancelled (exec.CommandContext kills it on ctx.Done), or when
// the caller invokes the returned kill function -- whichever comes first.
// Callers that want to force-stop a long-running stream (e.g. a monitor
// that outlives its usefulness) can call kill() instead of waiting for
// the context or the process to finish on its own.
//
// Both streams use a [util.BufferedPipe] so writers don't block indefinitely
// if the reader is slow. The buffer size is configured via
// [CliClientOptions.WithStreamBuffer]; nil means unbounded.
//
// kill is safe to call multiple times and after the process has already
// exited -- it closes the write ends of both pipes and sends SIGKILL to
// the subprocess. After kill returns, readers will drain any buffered
// data and then receive io.EOF.
//
// The returned streams emit io.EOF on clean exit (code 0), or an [ExitError]
// after the last byte has been consumed on non-zero exit.
func (c *cliClient) streamCmd(ctx context.Context, subcommand string, args ...string) (stdout, stderr io.Reader, kill func(), err error) {
	cmd := exec.CommandContext(ctx, c.binaryPath, args...)
	if len(c.env) > 0 {
		cmd.Env = append(os.Environ(), c.env...)
	}

	stdoutRead, stdoutWrite := util.NewBufferedPipe(util.BufferedPipeOptions{
		Ctx:    ctx,
		MaxBuf: c.maxBuf,
	})
	cmd.Stdout = stdoutWrite

	stderrRead, stderrWrite := util.NewBufferedPipe(util.BufferedPipeOptions{
		Ctx:    ctx,
		MaxBuf: c.maxBuf,
	})
	cmd.Stderr = stderrWrite

	if err := cmd.Start(); err != nil {
		_ = stdoutRead.Close()
		_ = stdoutWrite.Close()
		_ = stderrRead.Close()
		_ = stderrWrite.Close()
		return nil, nil, nil, newExecError(subcommand, c.binaryPath, args, "", "", err)
	}

	// kill closes both write ends and sends SIGKILL to the subprocess.
	// sync.Once makes it idempotent -- safe to call from the timeout
	// watcher, the exit goroutine, or the caller.
	var killOnce sync.Once
	killFn := func() {
		killOnce.Do(func() {
			_ = stdoutWrite.Close()
			_ = stderrWrite.Close()
			if cmd.Process != nil {
				_ = cmd.Process.Kill()
			}
		})
	}

	var setErr func(error)
	if c.timeout > 0 {
		stdout, setErr = util.NewTimeoutReader(stdoutRead, killFn, c.timeout)
	} else {
		stdout = stdoutRead
	}

	go func() {
		waitErr := cmd.Wait()
		if waitErr != nil {
			exitCode := 0
			var exitCodeErr *exec.ExitError
			if errors.As(waitErr, &exitCodeErr) {
				exitCode = exitCodeErr.ExitCode()
			}
			exitErr := newExitError(subcommand, c.binaryPath, args, exitCode, "", "", waitErr)
			killFn()
			if setErr != nil {
				setErr(exitErr)
			}
			return
		}
		killFn()
	}()

	return stdout, stderrRead, killFn, nil
}

// subcommandPath renders the leading verb tokens of an args slice
// for use as an error's subcommand. "cluster", "cluster show",
// "cluster update", etc. It stops at the first flag (token starting
// with '-') or at "--", whichever comes first.
func subcommandPath(args []string) string {
	var parts []string
	for _, a := range args {
		if a == "--" {
			break
		}
		if strings.HasPrefix(a, "-") {
			break
		}
		parts = append(parts, a)
	}
	return strings.Join(parts, " ")
}
