package version

import (
	"bytes"
	"errors"
	"io"
	"slices"
	"testing"
)

func TestGetLatestVersion(t *testing.T) {
	t.Parallel()

	got, ok := ParseList([]string{"v0.0.9", "v1.0.0", "v0.1.0", "release", "v01.0.0"}).Latest()
	if !ok || got != (Version{major: 1}) {
		t.Fatalf("Latest returned (%v, %t)", got, ok)
	}

	var isLatestVersion bool

	_, isLatestVersion = ParseList([]string{"release", "v1.0"}).Latest()
	if isLatestVersion {
		t.Fatal("Latest accepted invalid tags")
	}
}

func TestGetLatestVersionAcceptsZeroVersion(t *testing.T) {
	t.Parallel()

	got, found := ParseList([]string{"release", "v0.0.0", "v01.0.0"}).Latest()
	if !found || got != (Version{}) {
		t.Fatalf("Latest returned (%v, %t), want (v0.0.0, true)", got, found)
	}
}

func TestVersionListSortAndLast(t *testing.T) {
	t.Parallel()

	versions := ParseList([]string{"v1.10.0", "invalid", "v0.0.0", "v1.9.0", "v01.0.0"})
	before := slices.Clone(versions)

	latest, found := versions.Latest()
	if !found || latest != (Version{major: 1, minor: 10}) || !slices.Equal(versions, before) {
		t.Fatalf(
			"latest=%v found=%t versions=%v",
			latest,
			found,
			versions,
		)
	}

	versions.Sort()

	expected := VersionList{{}, {major: 1, minor: 9}, {major: 1, minor: 10}}
	if !slices.Equal(versions, expected) {
		t.Fatalf("sorted=%v want=%v", versions, expected)
	}

	for _, limit := range []int{0, 1, 2, 3, 10} {
		want := expected
		if limit > 0 && limit < len(expected) {
			want = expected[len(expected)-limit:]
		}

		got := versions.Last(limit)
		if !slices.Equal(got, want) {
			t.Fatalf(
				"Last(%d)=%v want=%v",
				limit,
				got,
				want,
			)
		}
	}
}

func TestVersionListEmpty(t *testing.T) {
	t.Parallel()

	versions := ParseList([]string{"release", "v01.0.0", "v1.0.0-beta"})
	if versions != nil {
		t.Fatalf("versions=%v", versions)
	}

	versions.Sort()

	latest, found := versions.Latest()
	if found || latest != (Version{}) || versions.Last(10) != nil {
		t.Fatalf("latest=%v found=%t", latest, found)
	}

	var stdout bytes.Buffer

	err := versions.Write(&stdout)
	if err != nil || stdout.Len() != 0 {
		t.Fatalf("error=%v output=%q", err, stdout.String())
	}
}

func TestVersionListWrite(t *testing.T) {
	t.Parallel()

	versions := VersionList{{}, {major: 1, minor: 10}}

	var stdout bytes.Buffer

	err := versions.Write(&stdout)
	if err != nil || stdout.String() != "v0.0.0\nv1.10.0\n" {
		t.Fatalf("error=%v output=%q", err, stdout.String())
	}

	err = versions.Write(failingVersionWriter{})
	if !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("error=%v", err)
	}
}

type failingVersionWriter struct{}

func (failingVersionWriter) Write(_ []byte) (int, error) { return 0, io.ErrClosedPipe }
