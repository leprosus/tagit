package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/leprosus/tagit/internal/cli"
)

func TestHelpWithoutGit(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	for _, argument := range []string{"help", "-h", "--help"} {
		var stdout, stderr bytes.Buffer

		code := cli.Run(
			t.Context(),
			[]string{argument},
			t.TempDir(),
			&stdout,
			&stderr,
		)
		if code != 0 || stdout.String() != cli.HelpMessage || stderr.Len() != 0 {
			t.Fatalf(
				"argument=%q code=%d stdout=%q stderr=%q",
				argument,
				code,
				stdout.String(),
				stderr.String(),
			)
		}
	}
}

func TestCommandHelpWithoutGit(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	for _, commandName := range []string{
		"patch", "minor", "major", "list", "latest", "show", "changes", "history", "set", "del", "del-all", "ver",
	} {
		for _, arguments := range [][]string{{"help", commandName}, {commandName, "-h"}, {commandName, "--help"}} {
			var stdout, stderr bytes.Buffer

			code := cli.Run(
				t.Context(),
				arguments,
				t.TempDir(),
				&stdout,
				&stderr,
			)
			if code != 0 || !strings.HasPrefix(stdout.String(), "Usage: tagit "+commandName) || stderr.Len() != 0 {
				t.Fatalf(
					"arguments=%q code=%d stdout=%q stderr=%q",
					arguments,
					code,
					stdout.String(),
					stderr.String(),
				)
			}
		}
	}
}

func TestHelpRejectsUnknownCommandAndExtraArguments(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	for _, arguments := range [][]string{{"help", "unknown"}, {"help", "list", "extra"}} {
		var stdout, stderr bytes.Buffer

		code := cli.Run(
			t.Context(),
			arguments,
			t.TempDir(),
			&stdout,
			&stderr,
		)
		if code != 1 || stdout.Len() != 0 || stderr.Len() == 0 {
			t.Fatalf(
				"arguments=%q code=%d stdout=%q stderr=%q",
				arguments,
				code,
				stdout.String(),
				stderr.String(),
			)
		}
	}
}
