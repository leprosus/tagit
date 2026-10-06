package command_test

import (
	"bytes"
	"testing"

	"github.com/leprosus/tagit/command"
	"github.com/leprosus/tagit/internal/cli"
)

func TestRunApplicationDeletesVersionTag(t *testing.T) {
	t.Parallel()

	directory := command.Repository(t)
	command.RunGit(
		t,
		directory,
		"tag",
		"v1.2.3",
	)

	var output bytes.Buffer

	err := cli.Execute(
		t.Context(),
		[]string{"del", "v1.2.3"},
		directory,
		&output,
	)
	if err != nil {
		t.Fatal(err)
	}

	if output.Len() != 0 {
		t.Fatalf("printed %q, want empty", output.String())
	}

	err = cli.Execute(
		t.Context(),
		[]string{"list"},
		directory,
		&output,
	)
	if err != nil {
		t.Fatal(err)
	}

	if output.Len() != 0 {
		t.Fatalf("list printed %q, want empty", output.String())
	}
}

func TestRunApplicationIgnoresMissingVersionTag(t *testing.T) {
	t.Parallel()

	directory := command.Repository(t)
	command.RunGit(
		t,
		directory,
		"tag",
		"v1.2.4",
	)

	var output bytes.Buffer

	err := cli.Execute(
		t.Context(),
		[]string{"del", "v1.2.3"},
		directory,
		&output,
	)
	if err != nil {
		t.Fatal(err)
	}

	if output.Len() != 0 {
		t.Fatalf("printed %q, want empty", output.String())
	}

	err = cli.Execute(
		t.Context(),
		[]string{"list"},
		directory,
		&output,
	)
	if err != nil {
		t.Fatal(err)
	}

	got := output.String()
	if got != "v1.2.4\n" {
		t.Fatalf("list printed %q, want v1.2.4\\n", got)
	}
}
