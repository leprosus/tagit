package git_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/leprosus/tagit/internal/git"
)

func TestCommitInfoWrite(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tag     string
		prefix  string
		message string
		want    string
	}{
		{name: "multiline", message: "Subject\nBody", want: "Subject"},
		{name: "CRLF", message: "Subject\r\nBody", want: "Subject"},
		{name: "empty", want: "(none)"},
		{name: "tagged", tag: "v1.2.3", prefix: "Tag: v1.2.3\n", message: "Subject", want: "Subject"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			info := git.CommitInfo{
				Tag:         test.tag,
				Commit:      "abcdef012345",
				ShortCommit: "abcdef0",
				Author:      "Alice <alice@example.com>",
				Date:        "2026-10-07T15:42:10+02:00",
				Message:     test.message,
			}

			want := test.prefix + "Commit: abcdef012345\nAuthor: Alice <alice@example.com>\n" +
				"Date: 2026-10-07T15:42:10+02:00\nMessage: " + test.want + "\n"

			var stdout bytes.Buffer

			err := info.Write(&stdout)
			if err != nil || stdout.String() != want {
				t.Fatalf(
					"error=%v got=%q want=%q",
					err,
					stdout.String(),
					want,
				)
			}

			stdout.Reset()

			err = info.WriteShort(&stdout)
			if err != nil || stdout.String() != "abcdef0: "+test.want+"\n" {
				t.Fatalf("error=%v output=%q", err, stdout.String())
			}
		})
	}
}

func TestCommitInfoListWrite(t *testing.T) {
	t.Parallel()

	commits := git.CommitInfoList{
		{Commit: "first", ShortCommit: "1111111", Author: "Alice", Date: "date1", Message: "First\nBody"},
		{Commit: "second", ShortCommit: "2222222", Author: "Bob", Date: "date2", Message: "Second"},
	}

	var stdout bytes.Buffer

	err := commits.Write(&stdout)

	want := "Commit: first\nAuthor: Alice\nDate: date1\nMessage: First\n\nCommit: second\nAuthor: Bob\nDate: date2\nMessage: Second\n"
	if err != nil || stdout.String() != want {
		t.Fatalf(
			"error=%v got=%q want=%q",
			err,
			stdout.String(),
			want,
		)
	}

	stdout.Reset()

	err = commits.WriteShort(&stdout)
	if err != nil || stdout.String() != "1111111: First\n2222222: Second\n" {
		t.Fatalf("error=%v output=%q", err, stdout.String())
	}
}

func TestCommitInfoListEmpty(t *testing.T) {
	t.Parallel()

	var (
		commits git.CommitInfoList
		stdout  bytes.Buffer
	)

	err := commits.Write(&stdout)
	if err != nil {
		t.Fatal(err)
	}

	err = commits.WriteShort(&stdout)
	if err != nil || stdout.Len() != 0 {
		t.Fatalf("error=%v output=%q", err, stdout.String())
	}
}

func TestCommitInfoWriteFailure(t *testing.T) {
	t.Parallel()

	info := git.CommitInfo{Tag: "v1.2.3"}

	commits := git.CommitInfoList{info}
	for _, write := range []func(io.Writer) error{info.Write, info.WriteShort, commits.Write, commits.WriteShort} {
		err := write(commitInfoFailingWriter{})
		if !errors.Is(err, io.ErrClosedPipe) {
			t.Fatalf("error=%v", err)
		}
	}
}

type commitInfoFailingWriter struct{}

func (commitInfoFailingWriter) Write(_ []byte) (int, error) { return 0, io.ErrClosedPipe }
