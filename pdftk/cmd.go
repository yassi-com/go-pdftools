package pdftk

import (
	"bytes"
	"context"
	"io"
	"os/exec"
)

// defaultExecutable is the binary invoked unless OptionExecutable names
// another one.
const defaultExecutable = "pdftk"

// command is the invocation an Option adjusts before it runs. Options mutate
// this spec rather than a constructed exec.Cmd, so replacing the executable
// does not have to rebuild a command that already carries argv and pipes.
type command struct {
	executable string
	args       []string
	stdin      io.Reader
	stdout     io.Writer
}

// newCommand builds an invocation of the default executable and applies
// options to it.
func newCommand(stdout io.Writer, stdin io.Reader, options []Option, args ...string) *command {
	cmd := &command{
		executable: defaultExecutable,
		args:       args,
		stdin:      stdin,
		stdout:     stdout,
	}

	for _, option := range options {
		option(cmd)
	}

	return cmd
}

// run executes the command, returning an *Error on failure.
func (c *command) run(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	var stderr bytes.Buffer

	cmd := exec.CommandContext(ctx, c.executable, c.args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = c.stdin, c.stdout, &stderr

	if err := cmd.Run(); err != nil {
		// A process killed by the context reports "signal: killed", which says
		// nothing about why. Report the context error instead so callers can
		// tell a cancellation from a genuine pdftk failure.
		if ctxErr := ctx.Err(); ctxErr != nil {
			err = ctxErr
		}

		return newError(c.executable, stderr.Bytes(), err)
	}

	return nil
}
