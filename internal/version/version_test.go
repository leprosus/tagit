package version

import (
	"errors"
	"math"
	"testing"

	errtypes "github.com/leprosus/tagit/internal/errors"
)

func TestParseVersion(t *testing.T) {
	t.Parallel()

	testCaseList := []struct {
		input string
		want  Version
		ok    bool
	}{
		{input: "v0.0.1", want: Version{patch: 1}, ok: true},
		{input: "v12.34.56", want: Version{major: 12, minor: 34, patch: 56}, ok: true},
		{input: "1.2.3", want: Version{}, ok: false},
		{input: "v01.2.3", want: Version{}, ok: false},
		{input: "v1.2", want: Version{}, ok: false},
		{input: "v1.2.3-beta", want: Version{}, ok: false},
		{input: "v18446744073709551616.0.0", want: Version{}, ok: false},
	}
	for _, test := range testCaseList {
		got, ok := Parse(test.input)
		if ok != test.ok {
			t.Fatalf("Parse(%q) is not success but %t is expected", test.input, test.ok)
		}

		if ok && got != test.want {
			t.Fatalf(
				"Parse(%q) = (%v, %t), want (%v, %t)",
				test.input,
				got,
				ok,
				test.want,
				test.ok,
			)
		}
	}
}

func TestVersionGetNextVersion(t *testing.T) {
	t.Parallel()

	base := Version{patch: 1}
	testCaseByKind := map[Kind]Version{
		Patch: {patch: 2},
		Minor: {minor: 1},
		Major: {major: 1, patch: 1},
	}

	for curKind, want := range testCaseByKind {
		got, err := base.Next(curKind)
		if err != nil || got != want {
			t.Errorf(
				"Next(%q) = (%v, %v), want (%v, nil)",
				curKind,
				got,
				err,
				want,
			)
		}
	}

	var err error

	_, err = base.Next("invalid")
	if err == nil {
		t.Fatal("Next(invalid) returned nil error")
	}

	_, err = (Version{patch: math.MaxUint64}).Next(Patch)

	var target *errtypes.VersionOverflowError
	if !errors.As(err, &target) || target.Kind != string(Patch) {
		t.Fatalf("overflow error = %v", err)
	}
}

func TestGetLatestVersion(t *testing.T) {
	t.Parallel()

	got, ok := Latest([]string{"v0.0.9", "v1.0.0", "v0.1.0", "release", "v01.0.0"})
	if !ok || got != (Version{major: 1}) {
		t.Fatalf("Latest returned (%v, %t)", got, ok)
	}

	var isLatestVersion bool

	_, isLatestVersion = Latest([]string{"release", "v1.0"})
	if isLatestVersion {
		t.Fatal("Latest accepted invalid tags")
	}
}

func TestGetLatestVersionAcceptsZeroVersion(t *testing.T) {
	t.Parallel()

	got, found := Latest([]string{"release", "v0.0.0", "v01.0.0"})
	if !found || got != (Version{}) {
		t.Fatalf("Latest returned (%v, %t), want (v0.0.0, true)", got, found)
	}
}
