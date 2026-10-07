package command_test

import (
	"bytes"
	"testing"

	"github.com/leprosus/tagit/command"
	"github.com/leprosus/tagit/internal/cli"
)

func TestCommitOutputIsConsistent(t *testing.T) {
	t.Parallel()
	directory := command.Repository(t)
	command.RunGit(
		t,
		directory,
		"tag",
		"v0.0.0",
	)
	hash := historyCommit(
		t,
		directory,
		"Subject\n\nCommit body excluded",
		"2026-10-07T15:42:10+02:00",
	)
	command.RunGit(
		t,
		directory,
		"tag",
		"-a",
		"v1.2.3",
		"-m",
		"Subject\n\nTag body excluded",
	)

	want := "Commit: " + hash + "\nAuthor: Alice <alice@example.com>\nDate: 2026-10-07T15:42:10+02:00\nMessage: Subject\n"

	for _, args := range [][]string{
		{"show", "v1.2.3"},
		{"changes", "v0.0.0"},
		{"history", "v0.0.0"},
	} {
		var stdout, stderr bytes.Buffer

		code := cli.Run(
			t.Context(),
			args,
			directory,
			&stdout,
			&stderr,
		)
		if code != 0 || stderr.Len() != 0 || stdout.String() != want {
			t.Fatalf(
				"args=%q code=%d stderr=%q got=%q want=%q",
				args,
				code,
				stderr.String(),
				stdout.String(),
				want,
			)
		}
	}
}
