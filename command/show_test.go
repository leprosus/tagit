package command_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/leprosus/tagit/command"
	"github.com/leprosus/tagit/internal/cli"
	"github.com/leprosus/tagit/internal/git"
)

func TestShow(t *testing.T) {
	t.Parallel()

	for _, annotated := range []bool{false, true} {
		name := "lightweight"
		if annotated {
			name = "annotated"
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assertShow(t, annotated)
		})
	}
}

func assertShow(t *testing.T, annotated bool) {
	t.Helper()
	directory := command.Repository(t)

	commit, err := git.New(directory).RunCommand(t.Context(), "rev-parse", "HEAD")
	if err != nil {
		t.Fatal(err)
	}

	message := "Message: initial\n"

	if annotated {
		command.RunGit(
			t,
			directory,
			"tag",
			"-a",
			"v1.2.3",
			"-m",
			"Release 1.2.3\n\nDetails",
		)
	} else {
		command.RunGit(
			t,
			directory,
			"tag",
			"v1.2.3",
		)
	}

	command.RunGit(
		t,
		directory,
		"commit",
		"--allow-empty",
		"-m",
		"Later commit",
	)
	before := repositoryTags(t, directory)

	var stdout, stderr bytes.Buffer

	code := cli.Run(
		t.Context(),
		[]string{"show", "v1.2.3"},
		directory,
		&stdout,
		&stderr,
	)
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("code=%d stderr=%q", code, stderr.String())
	}

	output := stdout.String()

	prefix := "Commit: " + strings.TrimSpace(commit) + "\nAuthor: Test User <test@example.com>\nDate: "
	if !strings.HasPrefix(output, prefix) || !strings.HasSuffix(output, message) {
		t.Fatalf("unexpected output: %q", output)
	}

	date, _, _ := strings.Cut(strings.TrimPrefix(output, prefix), "\n")

	_, err = time.Parse(time.RFC3339, date)
	if err != nil {
		t.Fatalf("invalid date %q: %v", date, err)
	}

	assertRepositoryTags(t, directory, before)
}

func TestShowMissingTag(t *testing.T) {
	t.Parallel()
	directory := command.Repository(t)
	command.RunGit(
		t,
		directory,
		"branch",
		"v1.2.3",
	)

	var stdout, stderr bytes.Buffer

	err := cli.Execute(
		t.Context(),
		[]string{"show", "v1.2.3"},
		directory,
		&stdout,
	)

	var target *git.TagNotFoundError
	if !errors.As(err, &target) || target.Tag != "v1.2.3" || stdout.Len() != 0 {
		t.Fatalf("error=%v stdout=%q", err, stdout.String())
	}

	code := cli.Run(
		t.Context(),
		[]string{"show", "v1.2.3"},
		directory,
		&stdout,
		&stderr,
	)
	if code != 1 || stdout.Len() != 0 || !strings.Contains(stderr.String(), `local tag "v1.2.3" not found`) {
		t.Fatalf(
			"code=%d stdout=%q stderr=%q",
			code,
			stdout.String(),
			stderr.String(),
		)
	}
}

func TestShowInvalidArguments(t *testing.T) {
	t.Parallel()
	directory := command.Repository(t)

	for _, args := range [][]string{{"show"}, {"show", "v01.2.3"}, {"show", "1.2.3"}, {"show", "v1.2.3-beta"}, {"show", "v1.2.3", "extra"}} {
		var stdout bytes.Buffer

		err := cli.Execute(
			t.Context(),
			args,
			directory,
			&stdout,
		)

		var target *command.UsageError
		if !errors.As(err, &target) || stdout.Len() != 0 {
			t.Fatalf(
				"args=%q error=%v stdout=%q",
				args,
				err,
				stdout.String(),
			)
		}
	}
}

func TestShowUsesCommitAuthor(t *testing.T) {
	t.Parallel()
	directory := command.Repository(t)
	historyCommit(
		t,
		directory,
		"Commit message\n\nBody excluded",
		"2026-10-07T15:42:10+02:00",
	)
	command.RunGit(
		t,
		directory,
		"tag",
		"-a",
		"v1.2.3",
		"-m",
		"Tag message",
	)

	var stdout bytes.Buffer

	err := command.Show(
		t.Context(),
		git.New(directory),
		[]string{"v1.2.3"},
		&stdout,
	)
	if err != nil {
		t.Fatal(err)
	}

	output := stdout.String()
	if !strings.Contains(output, "Author: Alice <alice@example.com>\nDate: 2026-10-07T15:42:10+02:00\n") ||
		!strings.HasSuffix(output, "Message: Commit message\n") {
		t.Fatalf("output=%q", output)
	}
}
