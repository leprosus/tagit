package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunApplicationShowsHelpOutsideRepository(t *testing.T) {
	t.Parallel()

	testCaseList := [][]string{
		nil,
		{"patch"},
		{"minor"},
		{"major"},
		{"list", "10"},
		{"set", "v1.2.3"},
		{"unknown"},
	}
	for _, argumentList := range testCaseList {
		var output bytes.Buffer

		err := runApplication(t.Context(), argumentList, t.TempDir(), &output)
		if err != nil {
			t.Fatalf("runApplication(%q): %v", argumentList, err)
		}

		got := output.String()
		if got != helpMessage {
			t.Errorf("runApplication(%q) printed %q, want %q", argumentList, got, helpMessage)
		}
	}
}

func TestRunApplicationReturnsUsageError(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer

	err := runApplication(t.Context(), []string{"patch", "extra"}, createRepository(t), &output)

	var target *usageError
	if !errors.As(err, &target) || !strings.Contains(err.Error(), "usage:") {
		t.Fatalf("error = %v", err)
	}
}

func TestRunCommandReturnsErrorExitCodeAndPrintsHelp(t *testing.T) {
	t.Parallel()

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	exitCode := runCommand(t.Context(), []string{"del", "1.2.3"}, createRepository(t), &stdout, &stderr)
	if exitCode != 1 {
		t.Fatalf("exit code = %d, want 1", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}

	got := stderr.String()

	if !strings.HasPrefix(got, "tagit: usage: tagit [patch|minor|major]\n") {
		t.Fatalf("stderr = %q, want error prefix", got)
	}

	if !strings.HasSuffix(got, helpMessage) {
		t.Errorf("stderr = %q, want help suffix %q", got, helpMessage)
	}
}

func TestRunApplicationReturnsCancelledContextError(t *testing.T) {
	t.Parallel()

	directory := createRepository(t)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	var output bytes.Buffer

	err := runApplication(ctx, []string{"list"}, directory, &output)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}

	if output.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", output.String())
	}
}

func TestRunApplicationReturnsMissingGitError(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	var output bytes.Buffer

	err := runApplication(t.Context(), []string{"list"}, t.TempDir(), &output)
	if !errors.Is(err, exec.ErrNotFound) {
		t.Fatalf("error = %v, want exec.ErrNotFound", err)
	}

	if output.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", output.String())
	}
}

func TestRunCommandReportsRepositoryCheckFailure(t *testing.T) {
	t.Parallel()

	directory := filepath.Join(t.TempDir(), "missing")

	var stdout, stderr bytes.Buffer

	exitCode := runCommand(t.Context(), []string{"list"}, directory, &stdout, &stderr)
	if exitCode != 1 || stdout.Len() != 0 || !strings.Contains(stderr.String(), "git rev-parse --git-dir:") {
		t.Fatalf("exit code = %d, stdout = %q, stderr = %q", exitCode, stdout.String(), stderr.String())
	}
}

func TestRunApplicationReturnsGitRepositoryAccessError(t *testing.T) {
	directory := t.TempDir()
	gitPath := filepath.Join(directory, "git")

	err := os.WriteFile(gitPath, []byte("#!/bin/sh\nprintf 'fatal: cannot access repository: Permission denied\\n' >&2\nexit 128\n"), 0o700)
	if err != nil {
		t.Fatal(err)
	}

	t.Setenv("PATH", directory)

	var output bytes.Buffer

	err = runApplication(t.Context(), []string{"list"}, directory, &output)

	var target *gitCommandError
	if !errors.As(err, &target) || !strings.Contains(err.Error(), "Permission denied") {
		t.Fatalf("error = %v, want Git access error", err)
	}

	if output.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", output.String())
	}
}

func createRepository(t *testing.T) (directory string) {
	t.Helper()
	directory = t.TempDir()
	runGitCommand(t, directory, "init")
	runGitCommand(t, directory, "config", "user.email", "test@example.com")
	runGitCommand(t, directory, "config", "user.name", "Test User")
	runGitCommand(t, directory, "commit", "--allow-empty", "-m", "initial")
	directory = filepath.Clean(directory)

	return directory
}

func runGitCommand(t *testing.T, directory string, argumentList ...string) {
	t.Helper()
	command := makeCommand(t.Context(), argumentList...)
	command.Dir = directory

	var (
		output []byte
		err    error
	)

	output, err = command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v: %s", strings.Join(argumentList, " "), err, output)
	}
}
