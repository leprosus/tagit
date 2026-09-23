package main

import (
	"bytes"
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunApplicationCreatesAndIncrementsTagList(t *testing.T) {
	t.Parallel()

	directory := createRepository(t)

	testCaseList := []struct {
		argument string
		want     string
	}{
		{"", "v0.0.1"},
		{"patch", "v0.0.2"},
		{"minor", "v0.1.0"},
		{"major", "v1.0.0"},
	}
	for _, test := range testCaseList {
		var (
			output       bytes.Buffer
			argumentList []string
		)
		if test.argument != "" {
			argumentList = []string{test.argument}
		}

		err := runApplication(argumentList, directory, &output)
		if err != nil {
			t.Fatalf("runApplication(%q): %v", test.argument, err)
		}

		got := strings.TrimSpace(output.String())
		if got != test.want {
			t.Errorf("runApplication(%q) printed %q, want %q", test.argument, got, test.want)
		}
	}
}

func TestRunApplicationUsesHighestSemanticVersionAndIgnoresOtherTagList(t *testing.T) {
	t.Parallel()

	directory := createRepository(t)
	for _, tag := range []string{"release", "v0.9.9", "v2.1.3", "v01.2.3"} {
		runGitCommand(t, directory, "tag", tag)
	}

	var output bytes.Buffer

	err := runApplication([]string{"patch"}, directory, &output)
	if err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(output.String())
	if got != "v2.1.4" {
		t.Fatalf("printed %q, want v2.1.4", got)
	}
}

func TestRunApplicationCreatesInitialTagWhenOnlyNonSemanticTagsExist(t *testing.T) {
	t.Parallel()

	directory := createRepository(t)
	runGitCommand(t, directory, "tag", "release-2026")

	var output bytes.Buffer

	err := runApplication([]string{"major"}, directory, &output)
	if err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(output.String())
	if got != "v0.0.1" {
		t.Fatalf("printed %q, want v0.0.1", got)
	}
}

func TestRunApplicationShowsHelpOutsideRepository(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer

	err := runApplication(nil, t.TempDir(), &output)
	if err != nil {
		t.Fatal(err)
	}

	got := output.String()
	if got != helpMessage {
		t.Fatalf("printed %q, want %q", got, helpMessage)
	}
}

func TestRunApplicationErrors(t *testing.T) {
	t.Parallel()

	t.Run("invalid command", func(t *testing.T) {
		t.Parallel()

		var output bytes.Buffer

		err := runApplication([]string{"build"}, createRepository(t), &output)

		var target *unknownVersionIncrementError
		if !errors.As(err, &target) || !strings.Contains(err.Error(), "unknown version increment") {
			t.Fatalf("error = %v", err)
		}
	})
	t.Run("too many arguments", func(t *testing.T) {
		t.Parallel()

		var output bytes.Buffer

		err := runApplication([]string{"patch", "extra"}, createRepository(t), &output)

		var target *usageError
		if !errors.As(err, &target) || !strings.Contains(err.Error(), "usage:") {
			t.Fatalf("error = %v", err)
		}
	})
	t.Run("not a repository", func(t *testing.T) {
		t.Parallel()

		var output bytes.Buffer

		err := runApplication([]string{"patch"}, t.TempDir(), &output)

		var target *gitCommandError
		if !errors.As(err, &target) || !strings.Contains(err.Error(), "git tag --list") {
			t.Fatalf("error = %v", err)
		}
	})
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
