package command_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/leprosus/tagit/command"
	"github.com/leprosus/tagit/internal/cli"
	"github.com/leprosus/tagit/internal/git"
)

func TestLatest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tags []string
		want string
	}{
		{name: "numeric", tags: []string{"v0.9.9", "v0.10.0"}, want: "v0.10.0\n"},
		{name: "greatest", tags: []string{"v3.1.2", "v2.99.99", "v3.1.3", "v3.0.9"}, want: "v3.1.3\n"},
		{name: "ignores invalid", tags: []string{"release", "v01.2.3", "v9.0.0-beta", "v1.2.3"}, want: "v1.2.3\n"},
		{name: "zero", tags: []string{"v0.0.0"}, want: "v0.0.0\n"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			directory := command.Repository(t)
			for _, tag := range test.tags {
				command.RunGit(
					t,
					directory,
					"tag",
					tag,
				)
			}

			before := repositoryTags(t, directory)

			var stdout, stderr bytes.Buffer

			code := cli.Run(
				t.Context(),
				[]string{"latest"},
				directory,
				&stdout,
				&stderr,
			)
			if code != 0 || stdout.String() != test.want || stderr.Len() != 0 {
				t.Fatalf(
					"code=%d stdout=%q stderr=%q want=%q",
					code,
					stdout.String(),
					stderr.String(),
					test.want,
				)
			}

			assertRepositoryTags(t, directory, before)
		})
	}
}

func TestLatestWithoutVersionTags(t *testing.T) {
	t.Parallel()

	for _, tags := range [][]string{nil, {"release", "v01.2.3", "v1.0.0-beta"}} {
		directory := command.Repository(t)
		for _, tag := range tags {
			command.RunGit(
				t,
				directory,
				"tag",
				tag,
			)
		}

		var stdout, stderr bytes.Buffer

		err := cli.Execute(
			t.Context(),
			[]string{"latest"},
			directory,
			&stdout,
		)

		var target *command.NoVersionTagsError
		if !errors.As(err, &target) || stdout.Len() != 0 {
			t.Fatalf("error=%v stdout=%q", err, stdout.String())
		}

		code := cli.Run(
			t.Context(),
			[]string{"latest"},
			directory,
			&stdout,
			&stderr,
		)
		if code != 1 || stdout.Len() != 0 || !strings.Contains(stderr.String(), "no local version tags found") {
			t.Fatalf(
				"code=%d stdout=%q stderr=%q",
				code,
				stdout.String(),
				stderr.String(),
			)
		}
	}
}

func TestLatestRejectsExtraArguments(t *testing.T) {
	t.Parallel()

	directory := command.Repository(t)
	command.RunGit(
		t,
		directory,
		"tag",
		"v1.2.3",
	)

	var stdout, stderr bytes.Buffer

	err := cli.Execute(
		t.Context(),
		[]string{"latest", "extra"},
		directory,
		&stdout,
	)

	if _, ok := errors.AsType[*command.UsageError](err); !ok {
		t.Fatalf("error=%v, want usage error", err)
	}

	code := cli.Run(
		t.Context(),
		[]string{"latest", "extra"},
		directory,
		&stdout,
		&stderr,
	)
	if code != 1 || stdout.Len() != 0 || stderr.Len() == 0 {
		t.Fatalf(
			"code=%d stdout=%q stderr=%q",
			code,
			stdout.String(),
			stderr.String(),
		)
	}

	command.AssertTagPresence(
		t,
		directory,
		"v1.2.3",
		true,
	)
}

func TestLatestOutsideRepository(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer

	code := cli.Run(
		t.Context(),
		[]string{"latest"},
		t.TempDir(),
		&stdout,
		&stderr,
	)
	if code != 0 || stdout.String() != cli.HelpMessage || stderr.Len() != 0 {
		t.Fatalf(
			"code=%d stdout=%q stderr=%q",
			code,
			stdout.String(),
			stderr.String(),
		)
	}
}

func repositoryTags(t *testing.T, directory string) (tags []string) {
	t.Helper()

	var err error

	tags, err = git.New(directory).GetTagList(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	return tags
}

func assertRepositoryTags(t *testing.T, directory string, before []string) {
	t.Helper()

	after := repositoryTags(t, directory)
	if strings.Join(before, "\n") != strings.Join(after, "\n") {
		t.Fatalf("tags changed: before=%q after=%q", before, after)
	}
}
