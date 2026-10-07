package command_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/leprosus/tagit/command"
	"github.com/leprosus/tagit/internal/cli"
	"github.com/leprosus/tagit/internal/git"
)

func TestChangesRanges(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tags []string
		args []string
		want []string
	}{
		{name: "explicit lightweight", tags: []string{"v1.2.3"}, args: []string{"v1.2.3"}, want: []string{"Second change", "First change"}},
		{name: "greatest version", tags: []string{"v1.9.0", "v1.10.0"}, want: []string{"Second change"}},
		{name: "entire history", want: []string{"Second change", "First change", "initial"}},
		{name: "ignores invalid tags", tags: []string{"release"}, want: []string{"Second change", "First change", "initial"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			directory := command.Repository(t)
			if len(test.tags) > 0 {
				command.RunGit(
					t,
					directory,
					"tag",
					test.tags[0],
				)
			}

			command.RunGit(
				t,
				directory,
				"commit",
				"--allow-empty",
				"-m",
				"First change\n\nBody excluded",
			)

			if len(test.tags) > 1 {
				command.RunGit(
					t,
					directory,
					"tag",
					test.tags[1],
				)
			}

			command.RunGit(
				t,
				directory,
				"commit",
				"--allow-empty",
				"-m",
				"Second change",
			)
			assertChanges(
				t,
				directory,
				test.args,
				test.want,
			)
		})
	}
}

func TestChangesAnnotatedAndEmpty(t *testing.T) {
	t.Parallel()
	directory := command.Repository(t)
	command.RunGit(
		t,
		directory,
		"tag",
		"-a",
		"v0.0.0",
		"-m",
		"Release",
	)
	assertChanges(
		t,
		directory,
		nil,
		nil,
	)
	command.RunGit(
		t,
		directory,
		"commit",
		"--allow-empty",
		"-m",
		"After annotated tag",
	)
	assertChanges(
		t,
		directory,
		[]string{"v0.0.0"},
		[]string{"After annotated tag"},
	)
}

func TestChangesOtherBranch(t *testing.T) {
	t.Parallel()
	directory := command.Repository(t)
	command.RunGit(
		t,
		directory,
		"checkout",
		"-b",
		"release",
	)
	command.RunGit(
		t,
		directory,
		"commit",
		"--allow-empty",
		"-m",
		"Release branch",
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
	command.RunGit(
		t,
		directory,
		"commit",
		"--allow-empty",
		"-m",
		"Work branch",
	)
	assertChanges(
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
	assertChanges(
		t,
		directory,
		[]string{"v2.0.0"},
		[]string{"Merge release", "Work branch"},
	)
}

func TestChangesErrors(t *testing.T) {
	t.Parallel()
	directory := command.Repository(t)
	command.RunGit(
		t,
		directory,
		"branch",
		"v1.2.3",
	)

	for _, args := range [][]string{
		{"changes", "v1.2.3"},
		{"changes", "v01.2.3"},
		{"changes", "1.2.3"},
		{"changes", "v1.2.3-beta"},
		{"changes", "v1.2.3", "extra"},
	} {
		var stdout, stderr bytes.Buffer

		code := cli.Run(
			t.Context(),
			args,
			directory,
			&stdout,
			&stderr,
		)
		if code != 1 || stdout.Len() != 0 || stderr.Len() == 0 {
			t.Fatalf(
				"args=%q code=%d stdout=%q stderr=%q",
				args,
				code,
				stdout.String(),
				stderr.String(),
			)
		}
	}

	var stdout bytes.Buffer

	err := command.Changes(
		t.Context(),
		git.New(directory),
		[]string{"v1.2.3"},
		&stdout,
	)

	if _, ok := errors.AsType[*git.TagNotFoundError](err); !ok {
		t.Fatalf("error=%v", err)
	}
}

func TestChangesGitFailure(t *testing.T) {
	t.Parallel()
	// A repository without commits has no HEAD to read.
	directory := t.TempDir()
	command.RunGit(t, directory, "init")

	var stdout bytes.Buffer

	err := command.Changes(
		t.Context(),
		git.New(directory),
		nil,
		&stdout,
	)

	var target *git.GitCommandError
	if !errors.As(err, &target) || stdout.Len() != 0 {
		t.Fatalf("error=%v stdout=%q", err, stdout.String())
	}
}

func assertChanges(t *testing.T, directory string, args, want []string) {
	t.Helper()
	before := repositoryTags(t, directory)

	var stdout, stderr bytes.Buffer

	code := cli.Run(
		t.Context(),
		append([]string{"changes"}, args...),
		directory,
		&stdout,
		&stderr,
	)
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("code=%d stderr=%q", code, stderr.String())
	}

	var messages []string

	for line := range strings.SplitSeq(strings.TrimSuffix(stdout.String(), "\n"), "\n") {
		if line == "" {
			continue
		}

		hash, message, found := strings.Cut(line, " ")
		if !found || len(hash) < 4 {
			t.Fatalf("invalid commit line %q", line)
		}

		messages = append(messages, message)
	}

	if strings.Join(messages, "\n") != strings.Join(want, "\n") {
		t.Fatalf("messages=%q want=%q", messages, want)
	}

	assertRepositoryTags(t, directory, before)
}

func TestChangesWriterFailure(t *testing.T) {
	t.Parallel()

	err := command.Changes(
		t.Context(),
		git.New(command.Repository(t)),
		nil,
		failingChangesWriter{},
	)
	if !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("error=%v", err)
	}
}

type failingChangesWriter struct{}

func (failingChangesWriter) Write(_ []byte) (int, error) { return 0, io.ErrClosedPipe }
