package command_test

import (
	"bytes"
	"testing"

	"github.com/leprosus/tagit/command"
	"github.com/leprosus/tagit/internal/cli"
)

func TestRunApplicationSetsVersionTag(t *testing.T) {
	t.Parallel()

	directory := command.Repository(t)

	var output bytes.Buffer

	err := cli.Execute(
		t.Context(),
		[]string{"set", "v1.2.3"},
		directory,
		&output,
	)
	if err != nil {
		t.Fatal(err)
	}

	got := output.String()
	if got != "v1.2.3\n" {
		t.Fatalf("printed %q, want v1.2.3\\n", got)
	}

	output.Reset()

	err = cli.Execute(
		t.Context(),
		[]string{"list"},
		directory,
		&output,
	)
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

	directory := command.Repository(t)

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

		err := cli.Execute(
			t.Context(),
			argumentList,
			directory,
			&output,
		)
		if err != nil {
			t.Fatalf("cli.Execute(%q): %v", argumentList, err)
		}

		got := output.String()
		if got != cli.HelpMessage {
			t.Errorf(
				"cli.Execute(%q) printed %q, want %q",
				argumentList,
				got,
				cli.HelpMessage,
			)
		}
	}
}
