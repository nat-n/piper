package piper_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp("", "piper-test-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create temp dir: %v\n", err)
		os.Exit(1)
	}

	binaryPath = filepath.Join(tmpDir, "pipedream")

	build := exec.Command("go", "build", "-o", binaryPath, "./testdata/pipedream")
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		os.RemoveAll(tmpDir)
		fmt.Fprintf(os.Stderr, "failed to build test fixture: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()
	os.RemoveAll(tmpDir)
	os.Exit(code)
}

func runPipedream(t *testing.T, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	cmd := exec.Command(binaryPath, args...)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Run()
	exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("failed to execute binary: %v", err)
		}
	}
	return stdoutBuf.String(), stderrBuf.String(), exitCode
}

func TestPipedream(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantExit   int
		wantStdout []string
		wantStderr []string
	}{
		{
			name:     "full_pipeline_verbose",
			args:     []string{"-v", "start", "Hello", "world!", "upper", "print", "3", "lower", "print", "1"},
			wantExit: 0,
			wantStdout: []string{
				" - creating some data for the pipeline with those two words",
				" - uppercasing those words",
				" - gonna print the words now",
				"HELLO WORLD!",
				" - lowercasing those words",
				"hello world!",
			},
		},
		{
			name:     "help_no_args",
			args:     []string{},
			wantExit: 0,
			wantStdout: []string{
				"* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *",
				"pipedream - dreamily pipes data through your tasks",
				"Usage:",
				"pipedream [global options] [command [arguments...] ...]",
				"Global options:",
				"-v  Verbose mode",
				"Commands:",
				"start - takes two words",
				"args: first word, second word",
				"upper - uppercase all the words",
				"lower - lowercase all the words",
				"print - print whatever is in the pipeline",
				"args: times",
			},
		},
		{
			name:     "error_unknown_flag",
			args:     []string{"-x"},
			wantExit: 0,
			wantStdout: []string{
				"Unknown flag: -x",
				"pipedream - dreamily pipes data through your tasks",
			},
		},
		{
			name:     "error_unknown_task",
			args:     []string{"notreal"},
			wantExit: 0,
			wantStdout: []string{
				"Unknown task: notreal",
				"pipedream - dreamily pipes data through your tasks",
			},
		},
		{
			name:     "error_insufficient_args",
			args:     []string{"start", "onlyoneword"},
			wantExit: 0,
			wantStdout: []string{
				"Insufficient arguments provided for task start, expected 2",
				"pipedream - dreamily pipes data through your tasks",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, exitCode := runPipedream(t, tt.args...)

			if exitCode != tt.wantExit {
				t.Errorf("exit code = %d, want %d\nstdout:\n%s\nstderr:\n%s",
					exitCode, tt.wantExit, stdout, stderr)
			}

			for _, want := range tt.wantStdout {
				if !strings.Contains(stdout, want) {
					t.Errorf("stdout missing substring %q\nfull stdout:\n%s", want, stdout)
				}
			}

			for _, want := range tt.wantStderr {
				if !strings.Contains(stderr, want) {
					t.Errorf("stderr missing substring %q\nfull stderr:\n%s", want, stderr)
				}
			}
		})
	}
}
