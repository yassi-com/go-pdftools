package pdftk

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// fakeExecutable writes an executable shell script and returns its path. It
// stands in for pdftk so these tests do not need the real binary installed.
func fakeExecutable(t *testing.T, body string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "fake-pdftk")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+body+"\n"), 0o700); err != nil {
		t.Fatalf("writing fake executable: %v", err)
	}

	return path
}

func TestNewCommand(t *testing.T) {
	tests := []struct {
		name           string
		options        []Option
		wantExecutable string
		wantArgs       []string
	}{
		{
			name:           "defaults to pdftk",
			wantExecutable: "pdftk",
			wantArgs:       []string{"input.pdf", "cat"},
		},
		{
			name:           "executable is replaced",
			options:        []Option{OptionExecutable("pdftk-java")},
			wantExecutable: "pdftk-java",
			wantArgs:       []string{"input.pdf", "cat"},
		},
		{
			name:           "flatten is appended to the args",
			options:        []Option{OptionFlatten()},
			wantExecutable: "pdftk",
			wantArgs:       []string{"input.pdf", "cat", "flatten"},
		},
		{
			name:           "options combine",
			options:        []Option{OptionExecutable("pdftk-java"), OptionFlatten()},
			wantExecutable: "pdftk-java",
			wantArgs:       []string{"input.pdf", "cat", "flatten"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newCommand(nil, nil, tt.options, "input.pdf", "cat")

			if cmd.executable != tt.wantExecutable {
				t.Errorf("executable = %q, want %q", cmd.executable, tt.wantExecutable)
			}
			if !reflect.DeepEqual(cmd.args, tt.wantArgs) {
				t.Errorf("args = %v, want %v", cmd.args, tt.wantArgs)
			}
		})
	}
}

func TestCommandRunAppliesOptionsToTheProcess(t *testing.T) {
	// Echo the arguments back so the assertion covers what actually reached
	// the process, not just the spec newCommand built.
	script := fakeExecutable(t, `printf '%s' "$*"`)

	var out bytes.Buffer
	cmd := newCommand(&out, nil, []Option{OptionExecutable(script), OptionFlatten()}, "input.pdf", "cat", "output", "-")

	if err := cmd.run(context.Background()); err != nil {
		t.Fatalf("run() = %v, want nil", err)
	}

	if got, want := out.String(), "input.pdf cat output - flatten"; got != want {
		t.Errorf("process args = %q, want %q", got, want)
	}
}

func TestCommandRunPipesStdinToStdout(t *testing.T) {
	script := fakeExecutable(t, `cat`)

	var out bytes.Buffer
	cmd := newCommand(&out, strings.NewReader("fdf data"), []Option{OptionExecutable(script)}, "-")

	if err := cmd.run(context.Background()); err != nil {
		t.Fatalf("run() = %v, want nil", err)
	}

	if got, want := out.String(), "fdf data"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
}

func TestCommandRunErrorKeepsStderrAndExitCode(t *testing.T) {
	script := fakeExecutable(t, `echo "OWNER PASSWORD REQUIRED" >&2; exit 3`)

	cmd := newCommand(nil, nil, []Option{OptionExecutable(script)}, "input.pdf")
	err := cmd.run(context.Background())

	var toolErr *Error
	if !errors.As(err, &toolErr) {
		t.Fatalf("run() = %v, want *Error", err)
	}

	if toolErr.Stderr != "OWNER PASSWORD REQUIRED" {
		t.Errorf("Stderr = %q, want %q", toolErr.Stderr, "OWNER PASSWORD REQUIRED")
	}
	if toolErr.ExitCode != 3 {
		t.Errorf("ExitCode = %d, want 3", toolErr.ExitCode)
	}
	if toolErr.Executable != script {
		t.Errorf("Executable = %q, want %q", toolErr.Executable, script)
	}

	// Callers classifying a failure need the stderr in the message and the
	// exec error still reachable underneath.
	if !strings.Contains(err.Error(), "OWNER PASSWORD REQUIRED") {
		t.Errorf("Error() = %q, want it to include stderr", err.Error())
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Errorf("errors.As(%v, *exec.ExitError) = false, want true", err)
	}
}

func TestCommandRunErrorOnMissingExecutable(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "not-installed")

	cmd := newCommand(nil, nil, []Option{OptionExecutable(missing)}, "input.pdf")
	err := cmd.run(context.Background())

	var toolErr *Error
	if !errors.As(err, &toolErr) {
		t.Fatalf("run() = %v, want *Error", err)
	}

	if toolErr.ExitCode != -1 {
		t.Errorf("ExitCode = %d, want -1", toolErr.ExitCode)
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("errors.Is(%v, os.ErrNotExist) = false, want true", err)
	}
}

func TestCommandRunReportsCanceledContext(t *testing.T) {
	script := fakeExecutable(t, `sleep 30`)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	cmd := newCommand(nil, nil, []Option{OptionExecutable(script)}, "input.pdf")

	if err := cmd.run(ctx); !errors.Is(err, context.Canceled) {
		t.Errorf("run() = %v, want context.Canceled", err)
	}
}

func TestCommandRunReportsExceededDeadline(t *testing.T) {
	script := fakeExecutable(t, `sleep 30`)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	cmd := newCommand(nil, nil, []Option{OptionExecutable(script)}, "input.pdf")

	// Killing the process reports "signal: killed" on its own; run must
	// surface the deadline instead.
	if err := cmd.run(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("run() = %v, want context.DeadlineExceeded", err)
	}
}
