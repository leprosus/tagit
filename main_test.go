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

		err := runApplication(t.Context(), argumentList, directory, &output)
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

	err := runApplication(t.Context(), []string{"patch"}, directory, &output)
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

	err := runApplication(t.Context(), []string{"major"}, directory, &output)
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

	testCaseList := [][]string{
		nil,
		{"patch"},
		{"minor"},
		{"major"},
		{"list", "10"},
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

func TestRunApplicationListsVersionTags(t *testing.T) {
	t.Parallel()

	directory := createRepository(t)
	for _, tag := range []string{"v1.0.0", "v0.10.0", "release", "v2.0.0", "v0.2.1", "v01.2.3"} {
		runGitCommand(t, directory, "tag", tag)
	}

	testCaseList := []struct {
		argumentList []string
		want         string
	}{
		{[]string{"list"}, "v0.2.1\nv0.10.0\nv1.0.0\nv2.0.0\n"},
		{[]string{"list", "2"}, "v1.0.0\nv2.0.0\n"},
	}
	for _, test := range testCaseList {
		var output bytes.Buffer

		err := runApplication(t.Context(), test.argumentList, directory, &output)
		if err != nil {
			t.Fatalf("runApplication(%q): %v", test.argumentList, err)
		}

		got := output.String()
		if got != test.want {
			t.Errorf("runApplication(%q) printed %q, want %q", test.argumentList, got, test.want)
		}
	}
}

func TestRunApplicationReturnsUnknownVersionIncrementError(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer

	err := runApplication(t.Context(), []string{"build"}, createRepository(t), &output)

	var target *unknownVersionIncrementError
	if !errors.As(err, &target) || !strings.Contains(err.Error(), "unknown version increment") {
		t.Fatalf("error = %v", err)
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

func TestRunApplicationReturnsInvalidListLimitError(t *testing.T) {
	t.Parallel()

	testCaseList := []string{"0", "many"}
	for _, value := range testCaseList {
		var output bytes.Buffer

		err := runApplication(t.Context(), []string{"list", value}, createRepository(t), &output)

		var target *invalidListLimitError
		if !errors.As(err, &target) || target.value != value {
			t.Fatalf("error = %v", err)
		}
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
