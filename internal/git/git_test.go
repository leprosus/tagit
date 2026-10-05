package git_test

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/leprosus/tagit/command"
	errtypes "github.com/leprosus/tagit/internal/errors"
	"github.com/leprosus/tagit/internal/git"
)

func TestIsRepository(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		directory string
		want      bool
	}{
		{command.Repository(t), true},
		{t.TempDir(), false},
	} {
		got, err := git.New(test.directory).IsRepository(t.Context())
		if err != nil || got != test.want {
			t.Fatalf("IsRepository(%q) = (%t, %v), want %t", test.directory, got, err, test.want)
		}
	}
}

func TestIsRepositoryReturnsCancelledContextError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	found, err := git.New(t.TempDir()).IsRepository(ctx)
	if found || !errors.Is(err, context.Canceled) {
		t.Fatalf("found=%t error=%v, want context.Canceled", found, err)
	}
}

func TestIsRepositoryReturnsMissingGitError(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	found, err := git.New(t.TempDir()).IsRepository(t.Context())
	if found || !errors.Is(err, exec.ErrNotFound) {
		t.Fatalf("found=%t error=%v, want exec.ErrNotFound", found, err)
	}
}

func TestIsRepositoryReturnsGitAccessError(t *testing.T) {
	directory := t.TempDir()
	gitPath := filepath.Join(directory, "git")

	err := os.WriteFile(gitPath, []byte("#!/bin/sh\nprintf 'fatal: cannot access repository: Permission denied\\n' >&2\nexit 128\n"), 0o700)
	if err != nil {
		t.Fatal(err)
	}

	t.Setenv("PATH", directory)
	found, err := git.New(directory).IsRepository(t.Context())

	var target *errtypes.GitCommandError
	if found || !errors.As(err, &target) || target.Output != "fatal: cannot access repository: Permission denied" {
		t.Fatalf("found=%t error=%v, want Git access error", found, err)
	}
}
