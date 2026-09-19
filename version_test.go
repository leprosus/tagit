package main

import (
	"errors"
	"math"
	"testing"
)

func TestParseVersion(t *testing.T) {
	t.Parallel()

	testCaseList := []struct {
		input string
		want  version
		ok    bool
	}{
		{"v0.0.1", version{0, 0, 1}, true},
		{"v12.34.56", version{12, 34, 56}, true},
		{"1.2.3", version{}, false},
		{"v01.2.3", version{}, false},
		{"v1.2", version{}, false},
		{"v1.2.3-beta", version{}, false},
		{"v18446744073709551616.0.0", version{}, false},
	}
	for _, test := range testCaseList {
		got, ok := parseVersion(test.input)
		if !ok && ok != test.ok {
			t.Fatalf("parseVersion(%q) is not success but %t is expected", test.input, test.ok)
		}

		if ok && got != test.want {
			t.Fatalf("parseVersion(%q) = (%v, %t), want (%v, %t)", test.input, got, ok, test.want, test.ok)
		}
	}
}

func TestVersionGetNextVersion(t *testing.T) {
	t.Parallel()

	base := version{0, 0, 1}
	testCaseByKind := map[kind]version{patch: {0, 0, 2}, minor: {0, 1, 0}, major: {1, 0, 1}}

	for curKind, want := range testCaseByKind {
		got, err := base.getNextVersion(curKind)
		if err != nil || got != want {
			t.Errorf("getNextVersion(%q) = (%v, %v), want (%v, nil)", curKind, got, err, want)
		}
	}

	var err error

	_, err = base.getNextVersion("invalid")
	if err == nil {
		t.Fatal("getNextVersion(invalid) returned nil error")
	}

	_, err = (version{patch: math.MaxUint64}).getNextVersion(patch)

	var target *versionOverflowError
	if !errors.As(err, &target) || target.kind != patch {
		t.Fatalf("overflow error = %v", err)
	}
}

func TestGetLatestVersion(t *testing.T) {
	t.Parallel()

	got, ok := getLatestVersion([]string{"v0.0.9", "v1.0.0", "v0.1.0", "release", "v01.0.0"})
	if !ok || got != (version{1, 0, 0}) {
		t.Fatalf("getLatestVersion returned (%v, %t)", got, ok)
	}

	var isLatestVersion bool

	_, isLatestVersion = getLatestVersion([]string{"release", "v1.0"})
	if isLatestVersion {
		t.Fatal("getLatestVersion accepted invalid tags")
	}
}
