package command_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/leprosus/tagit/command"
	"github.com/leprosus/tagit/internal/cli"
	"github.com/leprosus/tagit/internal/git"
)

func TestHistoryDetailedOutput(t *testing.T) {
	t.Parallel()
	directory := command.Repository(t)
	command.RunGit(
		t,
		directory,
		"tag",
		"v1.2.3",
	)
	first := historyCommit(
		t,
		directory,
		"First change\n\nDetails\n\nCommit: this is message text",
		"2026-10-07T14:18:03+02:00",
	)
	second := historyCommit(
		t,
		directory,
		"Second change",
		"2026-10-07T15:42:10+02:00",
	)
	want := "Commit: " + second + "\nAuthor: Alice <alice@example.com>\nDate: 2026-10-07T15:42:10+02:00\nMessage:\nSecond change\n\n" +
		"Commit: " + first + "\nAuthor: Alice <alice@example.com>\nDate: 2026-10-07T14:18:03+02:00\n" +
		"Message:\nFirst change\n\nDetails\n\nCommit: this is message text\n"
	before := repositoryTags(t, directory)

	var stdout, stderr bytes.Buffer

	code := cli.Run(
		t.Context(),
		[]string{"history", "v1.2.3"},
		directory,
		&stdout,
		&stderr,
	)
	if code != 0 || stderr.Len() != 0 || stdout.String() != want {
		t.Fatalf(
			"code=%d stderr=%q\ngot=%q\nwant=%q",
			code,
			stderr.String(),
			stdout.String(),
			want,
		)
	}

	assertRepositoryTags(t, directory, before)
}

func TestHistoryAutomaticRange(t *testing.T) {
	t.Parallel()
	directory := command.Repository(t)
	command.RunGit(
		t,
		directory,
		"tag",
		"v1.9.0",
	)
	historyCommit(
		t,
		directory,
		"Older change",
		"2026-10-07T14:18:03+02:00",
	)
	command.RunGit(
		t,
		directory,
		"tag",
		"-a",
		"v1.10.0",
		"-m",
		"Annotation excluded",
	)
	assertHistoryMessages(
		t,
		directory,
		nil,
		nil,
	)
	historyCommit(
		t,
		directory,
		"New change",
		"2026-10-07T15:42:10+02:00",
	)
	assertHistoryMessages(
		t,
		directory,
		nil,
		[]string{"New change"},
	)
	assertHistoryMessages(
		t,
		directory,
		[]string{"v1.10.0"},
		[]string{"New change"},
	)
}

func TestHistoryWithoutVersionTags(t *testing.T) {
	t.Parallel()
	directory := command.Repository(t)
	command.RunGit(
		t,
		directory,
		"tag",
		"release",
	)
	historyCommit(
		t,
		directory,
		"New change",
		"2026-10-07T15:42:10+02:00",
	)
	assertHistoryMessages(
		t,
		directory,
		nil,
		[]string{"New change", "initial"},
	)
}

func TestHistoryOtherBranchAndMerge(t *testing.T) {
	t.Parallel()
	directory := command.Repository(t)
	command.RunGit(
		t,
		directory,
		"checkout",
		"-b",
		"release",
	)
	historyCommit(
		t,
		directory,
		"Release branch",
		"2026-10-07T14:18:03+02:00",
	)
	command.RunGit(
		t,
		directory,
		"tag",
		"v2.0.0",
	)
	command.RunGit(
		t,
		directory,
		"checkout",
		"-b",
		"work",
		"HEAD~1",
	)
	historyCommit(
		t,
		directory,
		"Work branch",
		"2026-10-07T15:42:10+02:00",
	)
	assertHistoryMessages(
		t,
		directory,
		nil,
		[]string{"Work branch"},
	)
	command.RunGit(
		t,
		directory,
		"merge",
		"--no-ff",
		"release",
		"-m",
		"Merge release",
	)
	assertHistoryMessages(
		t,
		directory,
		[]string{"v2.0.0"},
		[]string{"Merge release", "Work branch"},
	)
}

func TestHistoryInvalidArguments(t *testing.T) {
	t.Parallel()
	curGit := git.New(command.Repository(t))

	for _, args := range [][]string{
		{"v01.2.3"}, {"1.2.3"}, {"v1.2.3-beta"}, {"v1.2.3", "extra"},
	} {
		var stdout bytes.Buffer

		err := command.History(
			t.Context(),
			curGit,
			args,
			&stdout,
		)
		if _, ok := errors.AsType[*command.UsageError](err); !ok || stdout.Len() != 0 {
			t.Fatalf(
				"args=%q error=%v stdout=%q",
				args,
				err,
				stdout.String(),
			)
		}
	}
}

func TestHistoryMissingTag(t *testing.T) {
	t.Parallel()
	directory := command.Repository(t)
	command.RunGit(
		t,
		directory,
		"branch",
		"v1.2.3",
	)

	var stdout, stderr bytes.Buffer

	err := command.History(
		t.Context(),
		git.New(directory),
		[]string{"v1.2.3"},
		&stdout,
	)
	if _, ok := errors.AsType[*git.TagNotFoundError](err); !ok {
		t.Fatalf("error=%v", err)
	}

	code := cli.Run(
		t.Context(),
		[]string{"history", "v1.2.3"},
		directory,
		&stdout,
		&stderr,
	)
	if code != 1 || stdout.Len() != 0 || !strings.Contains(stderr.String(), "local tag") {
		t.Fatalf(
			"code=%d stdout=%q stderr=%q",
			code,
			stdout.String(),
			stderr.String(),
		)
	}
}

func TestHistoryWriterFailure(t *testing.T) {
	t.Parallel()

	err := command.History(
		t.Context(),
		git.New(command.Repository(t)),
		nil,
		failingHistoryWriter{},
	)
	if !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("error=%v", err)
	}
}

type failingHistoryWriter struct{}

func (failingHistoryWriter) Write(_ []byte) (int, error) { return 0, io.ErrClosedPipe }

func historyCommit(t *testing.T, directory, message, date string) (hash string) {
	t.Helper()
	process := git.MakeCommand(
		t.Context(),
		"commit",
		"--allow-empty",
		"--author=Alice <alice@example.com>",
		"-m",
		message,
	)
	process.Dir = directory

	process.Env = append(os.Environ(), "GIT_AUTHOR_DATE=2020-01-01T00:00:00+00:00", "GIT_COMMITTER_DATE="+date)

	output, err := process.CombinedOutput()
	if err != nil {
		t.Fatalf("commit: %v: %s", err, output)
	}

	hash, err = git.New(directory).RunCommand(t.Context(), "rev-parse", "HEAD")
	if err != nil {
		t.Fatal(err)
	}

	return strings.TrimSpace(hash)
}

func assertHistoryMessages(t *testing.T, directory string, args, want []string) {
	t.Helper()

	var stdout, stderr bytes.Buffer

	code := cli.Run(
		t.Context(),
		append([]string{"history"}, args...),
		directory,
		&stdout,
		&stderr,
	)
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("code=%d stderr=%q", code, stderr.String())
	}

	output := stdout.String()
	if strings.Count(output, "Message:\n") != len(want) {
		t.Fatalf("output=%q want=%q", output, want)
	}

	for _, message := range want {
		_, rest, found := strings.Cut(output, "Message:\n"+message+"\n")
		if !found {
			t.Fatalf("missing or unordered message %q in %q", message, output)
		}

		output = rest
	}
}
