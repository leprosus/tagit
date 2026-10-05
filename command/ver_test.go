package command_test

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/leprosus/tagit/internal/cli"
)

func TestVerWorksWithoutGitOutsideRepository(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	var stdout, stderr bytes.Buffer

	code := cli.Run(t.Context(), []string{"ver"}, t.TempDir(), &stdout, &stderr)
	if code != 0 || stdout.String() != "dev\n" || stderr.Len() != 0 {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}

func TestVerRejectsExtraArguments(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	var stdout, stderr bytes.Buffer

	code := cli.Run(t.Context(), []string{"ver", "extra"}, t.TempDir(), &stdout, &stderr)
	if code != 1 || stdout.Len() != 0 || stderr.Len() == 0 {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}

func TestVerBuildVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		flags string
		want  string
	}{
		{"development", "", "dev\n"},
		{"release", "-s -w -X github.com/leprosus/tagit/command.releaseVersion=v1.2.3", "v1.2.3\n"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			binary := filepath.Join(t.TempDir(), "tagit")
			// All arguments come from fixed test fixtures and a temporary directory.
			build := exec.CommandContext(t.Context(), "go", "build", "-ldflags", test.flags, "-o", binary, "..") //nolint:gosec

			output, err := build.CombinedOutput()
			if err != nil {
				t.Fatalf("build: %v: %s", err, output)
			}

			invocation := exec.CommandContext(t.Context(), binary, "ver")
			invocation.Dir = t.TempDir()
			invocation.Env = []string{"PATH=" + t.TempDir()}

			output, err = invocation.CombinedOutput()
			if err != nil || string(output) != test.want {
				t.Fatalf("ver: output=%q error=%v want=%q", output, err, test.want)
			}
		})
	}
}
