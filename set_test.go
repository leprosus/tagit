package main

import (
	"bytes"
	"testing"
)

func TestRunApplicationSetsVersionTag(t *testing.T) {
	t.Parallel()

	directory := createRepository(t)

	var output bytes.Buffer

	err := runApplication(t.Context(), []string{"set", "v1.2.3"}, directory, &output)
	if err != nil {
		t.Fatal(err)
	}

	got := output.String()
	if got != "v1.2.3\n" {
		t.Fatalf("printed %q, want v1.2.3\\n", got)
	}

	output.Reset()

	err = runApplication(t.Context(), []string{"list"}, directory, &output)
	if err != nil {
		t.Fatal(err)
	}

	got = output.String()
	if got != "v1.2.3\n" {
		t.Fatalf("list printed %q, want v1.2.3\\n", got)
	}
}

func TestRunApplicationShowsHelpForInvalidSetVersionTag(t *testing.T) {
	t.Parallel()

	directory := createRepository(t)

	testCaseList := [][]string{
		{"set"},
		{"set", "1.2.3"},
		{"set", "v1.2"},
		{"set", "v01.2.3"},
		{"set", "v1.2.3-beta"},
		{"set", "v1.2.3", "extra"},
	}
	for _, argumentList := range testCaseList {
		var output bytes.Buffer

		err := runApplication(t.Context(), argumentList, directory, &output)
		if err != nil {
			t.Fatalf("runApplication(%q): %v", argumentList, err)
		}

		got := output.String()
		if got != helpMessage {
			t.Errorf("runApplication(%q) printed %q, want %q", argumentList, got, helpMessage)
		}
	}
}
