package pdftk

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// Error reports a failed invocation. It keeps stderr and the exit code so
// callers can classify the failure — an encrypted or corrupt input, say —
// instead of matching on a flattened message.
type Error struct {
	// Executable is the binary that was invoked.
	Executable string

	// Stderr is what the process wrote to stderr, trimmed. Empty when the
	// process never started.
	Stderr string

	// ExitCode is the process exit code, or -1 when it never started, was
	// signaled, or the context ended first.
	ExitCode int

	// Err is the underlying failure: an *exec.ExitError, a start-up error, or
	// the context error when the context ended before the process finished.
	Err error
}

func (e *Error) Error() string {
	if e.Stderr == "" {
		return fmt.Sprintf("%s failed: %v", e.Executable, e.Err)
	}

	return fmt.Sprintf("%s failed: %v: %s", e.Executable, e.Err, e.Stderr)
}

// Unwrap exposes the underlying error so errors.Is and errors.As reach
// *exec.ExitError and the context errors.
func (e *Error) Unwrap() error { return e.Err }

func newError(executable string, stderr []byte, err error) *Error {
	toolErr := &Error{
		Executable: executable,
		Stderr:     strings.TrimSpace(string(stderr)),
		ExitCode:   -1,
		Err:        err,
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		toolErr.ExitCode = exitErr.ExitCode()
	}

	return toolErr
}
