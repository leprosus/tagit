package main

import (
	"bytes"
	"errors"
	"testing"
)

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
