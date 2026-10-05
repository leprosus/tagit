package main

import (
	"bytes"
	"testing"
)

func TestRunApplicationDeletesVersionTag(t *testing.T) {
	t.Parallel()

	directory := createRepository(t)
	runGitCommand(t, directory, "tag", "v1.2.3")

	var output bytes.Buffer

	err := runApplication(t.Context(), []string{"del", "v1.2.3"}, directory, &output)
	if err != nil {
		t.Fatal(err)
	}

	if output.Len() != 0 {
		t.Fatalf("printed %q, want empty", output.String())
	}

	err = runApplication(t.Context(), []string{"list"}, directory, &output)
	if err != nil {
		t.Fatal(err)
	}

	if output.Len() != 0 {
		t.Fatalf("list printed %q, want empty", output.String())
	}
}

func TestRunApplicationIgnoresMissingVersionTag(t *testing.T) {
	t.Parallel()

	directory := createRepository(t)
	runGitCommand(t, directory, "tag", "v1.2.4")

	var output bytes.Buffer

	err := runApplication(t.Context(), []string{"del", "v1.2.3"}, directory, &output)
	if err != nil {
		t.Fatal(err)
	}

	if output.Len() != 0 {
		t.Fatalf("printed %q, want empty", output.String())
	}

	err = runApplication(t.Context(), []string{"list"}, directory, &output)
	if err != nil {
		t.Fatal(err)
	}

	got := output.String()
	if got != "v1.2.4\n" {
		t.Fatalf("list printed %q, want v1.2.4\\n", got)
	}
}
