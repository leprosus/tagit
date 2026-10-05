package command

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/leprosus/tagit/internal/git"
)

func Repository(t *testing.T) (directory string) {
	t.Helper()

	directory = t.TempDir()
	RunGit(t, directory, "init")
	RunGit(t, directory, "config", "user.email", "test@example.com")
	RunGit(t, directory, "config", "user.name", "Test User")
	RunGit(t, directory, "commit", "--allow-empty", "-m", "initial")
	directory = filepath.Clean(directory)

	return directory
}

func RunGit(t *testing.T, directory string, argumentList ...string) {
	t.Helper()

	command := git.MakeCommand(t.Context(), argumentList...)
	command.Dir = directory

	var (
		output []byte
		err    error
	)

	output, err = command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v: %s", strings.Join(argumentList, " "), err, output)
	}
}

func BareRemote(t *testing.T, directory, name string) string {
	t.Helper()

	remote := t.TempDir()
	RunGit(t, directory, "init", "--bare", remote)
	RunGit(t, directory, "remote", "add", name, remote)

	return remote
}

func AssertTagPresence(t *testing.T, directory, tag string, want bool) {
	t.Helper()

	found, err := git.New(directory).HasTag(t.Context(), tag)
	if err != nil || found != want {
		t.Fatalf("tag %s in %s: found=%t err=%v want=%t", tag, directory, found, err, want)
	}
}
