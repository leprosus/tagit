package main

import (
	"bytes"
	"errors"
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

func TestRunApplicationIncrementsZeroVersion(t *testing.T) {
	t.Parallel()

	testCaseList := []struct {
		argument string
		want     string
	}{
		{"", "v0.0.1"},
		{"patch", "v0.0.1"},
		{"minor", "v0.1.0"},
		{"major", "v1.0.0"},
	}
	for _, test := range testCaseList {
		directory := createRepository(t)
		runGitCommand(t, directory, "tag", "v0.0.0")

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

		if got := strings.TrimSpace(output.String()); got != test.want {
			t.Errorf("runApplication(%q) printed %q, want %q", test.argument, got, test.want)
		}

		runGitCommand(t, directory, "rev-parse", "--verify", "refs/tags/"+test.want)
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

func TestRunApplicationReturnsUnknownVersionIncrementError(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer

	err := runApplication(t.Context(), []string{"build"}, createRepository(t), &output)

	var target *unknownVersionIncrementError
	if !errors.As(err, &target) || !strings.Contains(err.Error(), "unknown version increment") {
		t.Fatalf("error = %v", err)
	}
}
