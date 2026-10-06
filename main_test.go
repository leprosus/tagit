package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestInterruptCancelsGitSubprocess(t *testing.T) {
	t.Parallel()

	binary := buildTestBinary(t)
	directory := filepath.Dir(binary)

	fakeGit := filepath.Join(directory, "git")

	err := os.WriteFile(fakeGit, []byte("#!/bin/sh\nprintf ready > \"$TEST_READY\"\nexec /bin/sleep 30\n"), 0o700)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	ready := filepath.Join(directory, "ready")
	process := exec.CommandContext(ctx, binary, "list")
	process.Dir = directory
	process.Env = []string{"PATH=" + directory, "TEST_READY=" + ready}

	var stdout, stderr bytes.Buffer

	process.Stdout = &stdout
	process.Stderr = &stderr

	err = process.Start()
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan error, 1)
	go func() { done <- process.Wait() }()

	waitForGit(
		t,
		ctx,
		ready,
		done,
	)

	err = process.Process.Signal(os.Interrupt)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case err = <-done:
		failedAsExpected := err != nil && process.ProcessState.ExitCode() == 1
		outputIsCorrect := stdout.Len() == 0 && strings.Contains(stderr.String(), "git rev-parse")

		if !failedAsExpected || !outputIsCorrect {
			t.Fatalf(
				"error=%v code=%d stdout=%q stderr=%q",
				err,
				process.ProcessState.ExitCode(),
				stdout.String(),
				stderr.String(),
			)
		}
	case <-ctx.Done():
		t.Fatal("interrupt did not stop the Git subprocess")
	}
}

func buildTestBinary(t *testing.T) (binary string) {
	t.Helper()
	directory := t.TempDir()
	binary = filepath.Join(directory, "tagit")
	build := exec.CommandContext(
		t.Context(),
		"go",
		"build",
		"-o",
		binary,
		".",
	)

	output, err := build.CombinedOutput()
	if err != nil {
		t.Fatalf("build: %v: %s", err, output)
	}

	return binary
}

func waitForGit(
	t *testing.T,
	ctx context.Context,
	ready string,
	done <-chan error,
) {
	t.Helper()

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		_, err := os.Stat(ready)
		if err == nil {
			break
		}

		select {
		case err := <-done:
			t.Fatalf("process exited before Git started: %v", err)
		case <-ctx.Done():
			t.Fatal("Git subprocess did not start")
		case <-ticker.C:
		}
	}
}
