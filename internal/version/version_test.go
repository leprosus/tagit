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
		{"v0.0.1", Version{0, 0, 1}, true},
		{"v12.34.56", Version{12, 34, 56}, true},
		{"1.2.3", Version{}, false},
		{"v01.2.3", Version{}, false},
		{"v1.2", Version{}, false},
		{"v1.2.3-beta", Version{}, false},
		{"v18446744073709551616.0.0", Version{}, false},
	}
	for _, test := range testCaseList {
		got, ok := Parse(test.input)
		if !ok && ok != test.ok {
			t.Fatalf("Parse(%q) is not success but %t is expected", test.input, test.ok)
		}

		if ok && got != test.want {
			t.Fatalf("Parse(%q) = (%v, %t), want (%v, %t)", test.input, got, ok, test.want, test.ok)
		}
	}
}

func TestVersionGetNextVersion(t *testing.T) {
	t.Parallel()

	base := Version{0, 0, 1}
	testCaseByKind := map[Kind]Version{Patch: {0, 0, 2}, Minor: {0, 1, 0}, Major: {1, 0, 1}}

	for curKind, want := range testCaseByKind {
		got, err := base.Next(curKind)
		if err != nil || got != want {
			t.Errorf("Next(%q) = (%v, %v), want (%v, nil)", curKind, got, err, want)
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
	if !ok || got != (Version{1, 0, 0}) {
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
