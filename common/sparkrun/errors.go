package sparkrun

import (
	"errors"
	"fmt"

	"github.com/Riven-Spell/generic/ptr_tools"
	"github.com/Riven-Spell/sparky/common/cmdline"
	"github.com/Riven-Spell/sparky/common/enum"
)

// SparkrunError is implemented by every error this package returns.
//
// Callers that only need the command output can read the common fields
// (Binary, Args, StdError, StdOut, ExitCode) directly. Callers that
// need per-kind structure should type-switch on the concrete kinds:
// [ParseError], [ExecError], [ExitError], [UsageError], [TargetError];
// or compare Kind() against enum.ESparkrunErrorKind.
//
// The trailing implSparkrunError method seals the interface so only
// this package can satisfy it.
type SparkrunError interface {
	error

	Kind() enum.SparkrunErrorKind
	Binary() string
	Args(cmdline.Tag) []string
	StdError() string
	StdOut() string
	ExitCode() (valid bool, code int)

	implSparkrunError()
}

// sparkrunError is the shared value-carrying base for every kind of
// [SparkrunError].
//
//   - subcommand is the full nested path (e.g. "cluster show").
//   - err is the underlying cause; message rendering always begins
//     with the literal "sparkrun".
//   - kind classifies the error one of enum.SparkrunErrorKind values.
//   - binary and args identify the sparkrun invocation.
//   - stdErr/stdOut/exitCode hold the command output captured from
//     sparkrun on a non-zero exit.
type sparkrunError struct {
	subcommand string
	err        error

	kind   enum.SparkrunErrorKind
	binary string
	args   []string

	stdErr   string
	stdOut   string
	exitCode *int
}

func (s sparkrunError) Error() string {
	if s.subcommand == "" {
		return fmt.Sprintf("sparkrun: %v", s.err)
	}
	return fmt.Sprintf("sparkrun %s: %v", s.subcommand, s.err)
}

func (s sparkrunError) StdError() string {
	return s.stdErr
}

func (s sparkrunError) Kind() enum.SparkrunErrorKind {
	return s.kind
}

func (s sparkrunError) StdOut() string {
	return s.stdOut
}

func (s sparkrunError) Binary() string {
	return s.binary
}

func (s sparkrunError) Args(cmdline.Tag) []string {
	return s.args
}

func (s sparkrunError) ExitCode() (bool, int) {
	return s.exitCode != nil, ptr_tools.DerefOrZero(s.exitCode)
}

func (s sparkrunError) implSparkrunError() {}

// Unwrap exposes the underlying cause so errors.Is / errors.As can
// reach it.
func (s sparkrunError) Unwrap() error {
	return s.err
}

// AsSparkrunError unwraps err to the common [SparkrunError], if it is
// (or wraps) one. It returns (nil, false) for any other error.
func AsSparkrunError(err error) (SparkrunError, bool) {
	var target SparkrunError
	if errors.As(err, &target) {
		return target, true
	}
	return nil, false
}

// ParseError indicates sparkrun exited cleanly (exit 0) but its
// --json body was malformed. raw holds the undecodable stdout.
type ParseError struct {
	sparkrunError
	raw []byte
}

// Raw returns the undecodable --json body.
func (e *ParseError) Raw() []byte {
	if e == nil {
		return nil
	}
	return e.raw
}

func newParseError(subcommand, binary string, args []string, raw []byte, stdErr string, cause error) error {
	return &ParseError{
		sparkrunError: sparkrunError{
			subcommand: subcommand,
			err:        fmt.Errorf("%w: --json output did not decode", cause),
			kind:       enum.ESparkrunErrorKind.Parse(),
			binary:     binary,
			args:       append([]string(nil), args...),
			stdErr:     stdErr,
		},
		raw: append([]byte(nil), raw...),
	}
}

// ExecError indicates the sparkrun subprocess could not be started.
// Binary and Args identify the failed invocation (via the shared
// sparkrunError base).
type ExecError struct {
	sparkrunError
}

func newExecError(subcommand, binary string, args []string, stdErr, stdOut string, cause error) error {
	return &ExecError{
		sparkrunError: sparkrunError{
			subcommand: subcommand,
			err:        fmt.Errorf("could not start subprocess: %w", cause),
			kind:       enum.ESparkrunErrorKind.Exec(),
			binary:     binary,
			args:       append([]string(nil), args...),
			stdErr:     stdErr,
			stdOut:     stdOut,
		},
	}
}

// ExitError indicates the sparkrun subprocess exited with a non-zero
// code. exitCode, StdError, and StdOut describe the failed run.
type ExitError struct {
	sparkrunError
	result *ClusterCheckJobResult
	raw    []byte
}

// Result returns the decoded body, if any, that sparkrun emitted
// alongside a non-zero exit (set by [Client.ClusterCheckJob]).
func (e *ExitError) Result() *ClusterCheckJobResult {
	if e == nil {
		return nil
	}
	return e.result
}

// Raw returns the undecodable stdout body, if any, that sparkrun
// emitted alongside a non-zero exit.
func (e *ExitError) Raw() []byte {
	if e == nil {
		return nil
	}
	return e.raw
}

func newExitError(subcommand, binary string, args []string, exitCode int, stdErr, stdOut string, cause error) error {
	return &ExitError{
		sparkrunError: sparkrunError{
			subcommand: subcommand,
			err:        fmt.Errorf("exit code %d: %w", exitCode, cause),
			kind:       enum.ESparkrunErrorKind.Exit(),
			binary:     binary,
			args:       append([]string(nil), args...),
			stdErr:     stdErr,
			stdOut:     stdOut,
			exitCode:   &exitCode,
		},
	}
}

// UsageError indicates sparkrun rejected the invocation's arguments
// (click's exit-2 invalid-flag case). It is a non-zero exit, but
// classified separately so callers can distinguish caller error from
// sparkrun failure.
type UsageError struct {
	sparkrunError
}

func newUsageError(subcommand, binary string, args []string, exitCode int, stdErr, stdOut string, cause error) error {
	return &UsageError{
		sparkrunError: sparkrunError{
			subcommand: subcommand,
			err:        fmt.Errorf("usage: %w", cause),
			kind:       enum.ESparkrunErrorKind.Usage(),
			binary:     binary,
			args:       append([]string(nil), args...),
			stdErr:     stdErr,
			stdOut:     stdOut,
			exitCode:   &exitCode,
		},
	}
}

// TargetError indicates a method was called with a target kind it does
// not accept (e.g. a recipe file passed to Stop).
type TargetError struct {
	sparkrunError
	expected []string
	got      string
}

// Expected returns the target kinds the method accepts.
func (e *TargetError) Expected() []string {
	if e == nil {
		return nil
	}
	return e.expected
}

// Got returns the target kind that was actually passed.
func (e *TargetError) Got() string {
	if e == nil {
		return ""
	}
	return e.got
}

func newTargetError(subcommand, binary string, args []string, expected []string, got string) error {
	return &TargetError{
		sparkrunError: sparkrunError{
			subcommand: subcommand,
			err:        fmt.Errorf("target got %q, expected one of %v", got, expected),
			kind:       enum.ESparkrunErrorKind.Target(),
			binary:     binary,
			args:       append([]string(nil), args...),
		},
		expected: append([]string(nil), expected...),
		got:      got,
	}
}
