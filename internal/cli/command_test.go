package cli_test

import (
	"bytes"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/leprosus/tagit/command"
	"github.com/leprosus/tagit/internal/cli"

	errtypes "github.com/leprosus/tagit/internal/errors"
)

func TestRunApplicationShowsHelpOutsideRepository(t *testing.T) {
	t.Parallel()

	testCaseList := [][]string{
		nil,
		{"patch"},
		{"minor"},
		{"major"},
		{"list", "10"},
		{"set", "v1.2.3"},
		{"unknown"},
	}
	for _, argumentList := range testCaseList {
		var output bytes.Buffer

		err := cli.Execute(
			t.Context(),
			argumentList,
			t.TempDir(),
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

func TestRunApplicationReturnsUsageError(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer

	err := cli.Execute(
		t.Context(),
		[]string{"patch", "extra"},
		command.Repository(t),
		&output,
	)

	var target *errtypes.UsageError
	if !errors.As(err, &target) || !strings.Contains(err.Error(), "usage:") {
		t.Fatalf("error = %v", err)
	}
}

func TestRunCommandReturnsErrorExitCodeAndPrintsHelp(t *testing.T) {
	t.Parallel()

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	exitCode := cli.Run(
		t.Context(),
		[]string{"del", "1.2.3"},
		command.Repository(t),
		&stdout,
		&stderr,
	)
	if exitCode != 1 {
		t.Fatalf("exit code = %d, want 1", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}

	got := stderr.String()

	if !strings.HasPrefix(got, "tagit: usage: tagit [patch|minor|major]\n") {
		t.Fatalf("stderr = %q, want error prefix", got)
	}

	if !strings.HasSuffix(got, cli.HelpMessage) {
		t.Errorf("stderr = %q, want help suffix %q", got, cli.HelpMessage)
	}
}

func TestRunCommandReportsRepositoryCheckFailure(t *testing.T) {
	t.Parallel()

	directory := filepath.Join(t.TempDir(), "missing")

	var stdout, stderr bytes.Buffer

	exitCode := cli.Run(
		t.Context(),
		[]string{"list"},
		directory,
		&stdout,
		&stderr,
	)
	if exitCode != 1 || stdout.Len() != 0 || !strings.Contains(stderr.String(), "git rev-parse --git-dir:") {
		t.Fatalf(
			"exit code = %d, stdout = %q, stderr = %q",
			exitCode,
			stdout.String(),
			stderr.String(),
		)
	}
}
