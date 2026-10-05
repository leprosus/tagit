package command_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/leprosus/tagit/command"
	"github.com/leprosus/tagit/internal/cli"
	errtypes "github.com/leprosus/tagit/internal/errors"
)

func TestRunApplicationListsVersionTags(t *testing.T) {
	t.Parallel()

	directory := command.Repository(t)
	for _, tag := range []string{"v1.0.0", "v0.10.0", "release", "v2.0.0", "v0.2.1", "v01.2.3"} {
		command.RunGit(t, directory, "tag", tag)
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

		err := cli.Execute(t.Context(), test.argumentList, directory, &output)
		if err != nil {
			t.Fatalf("cli.Execute(%q): %v", test.argumentList, err)
		}

		got := output.String()
		if got != test.want {
			t.Errorf("cli.Execute(%q) printed %q, want %q", test.argumentList, got, test.want)
		}
	}
}

func TestRunApplicationReturnsInvalidListLimitError(t *testing.T) {
	t.Parallel()

	testCaseList := []string{"0", "many"}
	for _, value := range testCaseList {
		var output bytes.Buffer

		err := cli.Execute(t.Context(), []string{"list", value}, command.Repository(t), &output)

		var target *errtypes.InvalidListLimitError
		if !errors.As(err, &target) || target.Value != value {
			t.Fatalf("error = %v", err)
		}
	}
}
