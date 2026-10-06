package command_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/leprosus/tagit/internal/git"

	"github.com/leprosus/tagit/command"
	"github.com/leprosus/tagit/internal/cli"
)

func TestDeleteAllTags(t *testing.T) {
	t.Parallel()
	directory := command.Repository(t)
	command.RunGit(
		t,
		directory,
		"tag",
		"v1.2.3",
	)
	command.RunGit(
		t,
		directory,
		"tag",
		"v1.2.4",
	)

	remotes := []string{command.BareRemote(t, directory, "origin"), command.BareRemote(t, directory, "upstream")}
	for _, remote := range remotes {
		command.RunGit(
			t,
			directory,
			"push",
			remote,
			"refs/tags/v1.2.3",
			"refs/tags/v1.2.4",
		)
	}

	var output bytes.Buffer

	err := cli.Execute(
		t.Context(),
		[]string{"del-all", "v1.2.3"},
		directory,
		&output,
	)
	if err != nil {
		t.Fatal(err)
	}

	for _, repository := range append(remotes, directory) {
		command.AssertTagPresence(
			t,
			repository,
			"v1.2.3",
			false,
		)
		command.AssertTagPresence(
			t,
			repository,
			"v1.2.4",
			true,
		)
	}

	for _, name := range []string{"origin:", "upstream:", "local:"} {
		if !strings.Contains(output.String(), name) {
			t.Errorf("output = %q, missing %q", output.String(), name)
		}
	}

	err = cli.Execute(
		t.Context(),
		[]string{"del-all", "v1.2.3"},
		directory,
		&output,
	)
	if err != nil {
		t.Fatalf("retry: %v", err)
	}
}

func TestDeleteAllTagsPartialFailureAndRetry(t *testing.T) {
	t.Parallel()
	directory := command.Repository(t)
	command.RunGit(
		t,
		directory,
		"tag",
		"v1.2.3",
	)
	missing := filepath.Join(t.TempDir(), "missing.git")
	command.RunGit(
		t,
		directory,
		"remote",
		"add",
		"a-broken",
		missing,
	)
	remote := command.BareRemote(t, directory, "z-working")
	command.RunGit(
		t,
		directory,
		"push",
		remote,
		"refs/tags/v1.2.3",
	)

	var stdout, stderr bytes.Buffer

	code := cli.Run(
		t.Context(),
		[]string{"del-all", "v1.2.3"},
		directory,
		&stdout,
		&stderr,
	)
	if code != 1 || !strings.Contains(stderr.String(), "a-broken") || !strings.Contains(stdout.String(), "z-working:") {
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
	command.AssertTagPresence(
		t,
		remote,
		"v1.2.3",
		false,
	)
	command.RunGit(
		t,
		directory,
		"init",
		"--bare",
		missing,
	)
	command.RunGit(
		t,
		directory,
		"push",
		missing,
		"refs/tags/v1.2.3",
	)

	err := cli.Execute(
		t.Context(),
		[]string{"del-all", "v1.2.3"},
		directory,
		&stdout,
	)
	if err != nil {
		t.Fatal(err)
	}

	command.AssertTagPresence(
		t,
		directory,
		"v1.2.3",
		false,
	)
	command.AssertTagPresence(
		t,
		missing,
		"v1.2.3",
		false,
	)
}

func TestDeleteAllTagsWithoutRemotes(t *testing.T) {
	t.Parallel()
	directory := command.Repository(t)
	command.RunGit(
		t,
		directory,
		"tag",
		"v1.2.3",
	)

	var output bytes.Buffer
	for range 2 {
		err := cli.Execute(
			t.Context(),
			[]string{"del-all", "v1.2.3"},
			directory,
			&output,
		)
		if err != nil {
			t.Fatal(err)
		}

		command.AssertTagPresence(
			t,
			directory,
			"v1.2.3",
			false,
		)
	}
}

func TestDeleteAllTagsRejectsInvalidArguments(t *testing.T) {
	t.Parallel()
	directory := command.Repository(t)
	command.RunGit(
		t,
		directory,
		"tag",
		"v1.2.3",
	)

	for _, args := range [][]string{
		{"del-all"},
		{"del-all", "1.2.3"},
		{"del-all", "v01.2.3"},
		{"del-all", "v1.2.3", "extra"},
	} {
		var stdout, stderr bytes.Buffer

		code := cli.Run(
			t.Context(),
			args,
			directory,
			&stdout,
			&stderr,
		)
		if code != 1 || stdout.Len() != 0 {
			t.Fatalf(
				"args=%q code=%d stdout=%q",
				args,
				code,
				stdout.String(),
			)
		}

		command.AssertTagPresence(
			t,
			directory,
			"v1.2.3",
			true,
		)
	}
}

func TestDeleteAllTagsUsesAllPushURLs(t *testing.T) {
	t.Parallel()
	directory := command.Repository(t)
	command.RunGit(
		t,
		directory,
		"tag",
		"v1.2.3",
	)
	fetch := command.BareRemote(t, directory, "origin")
	first := t.TempDir()

	second := t.TempDir()
	for _, remote := range []string{fetch, first, second} {
		command.RunGit(
			t,
			directory,
			"init",
			"--bare",
			remote,
		)
		command.RunGit(
			t,
			directory,
			"push",
			remote,
			"refs/tags/v1.2.3",
		)
	}

	for _, remote := range []string{first, second} {
		command.RunGit(
			t,
			directory,
			"remote",
			"set-url",
			"--add",
			"--push",
			"origin",
			remote,
		)
	}

	var output bytes.Buffer

	err := cli.Execute(
		t.Context(),
		[]string{"del-all", "v1.2.3"},
		directory,
		&output,
	)
	if err != nil {
		t.Fatal(err)
	}

	for _, repository := range []string{first, second, directory} {
		command.AssertTagPresence(
			t,
			repository,
			"v1.2.3",
			false,
		)
	}

	command.AssertTagPresence(
		t,
		fetch,
		"v1.2.3",
		true,
	)
}

func TestDeleteAllTagsHandlesRemoteOnlyTag(t *testing.T) {
	t.Parallel()
	directory := command.Repository(t)
	command.RunGit(
		t,
		directory,
		"tag",
		"v1.2.3",
	)
	remote := command.BareRemote(t, directory, "origin")
	command.RunGit(
		t,
		directory,
		"push",
		remote,
		"refs/tags/v1.2.3",
	)
	command.RunGit(
		t,
		directory,
		"tag",
		"--delete",
		"v1.2.3",
	)

	var output bytes.Buffer

	err := cli.Execute(
		t.Context(),
		[]string{"del-all", "v1.2.3"},
		directory,
		&output,
	)
	if err != nil {
		t.Fatal(err)
	}

	command.AssertTagPresence(
		t,
		remote,
		"v1.2.3",
		false,
	)
}

func TestDeleteAllTagsHandlesRejectedPush(t *testing.T) {
	t.Parallel()
	directory := command.Repository(t)
	command.RunGit(
		t,
		directory,
		"tag",
		"v1.2.3",
	)
	rejected := command.BareRemote(t, directory, "a-rejected")

	working := command.BareRemote(t, directory, "z-working")
	for _, remote := range []string{rejected, working} {
		command.RunGit(
			t,
			directory,
			"push",
			remote,
			"refs/tags/v1.2.3",
		)
	}

	err := os.WriteFile(filepath.Join(rejected, "hooks", "pre-receive"), []byte("#!/bin/sh\nexit 1\n"), 0o700)
	if err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer

	err = cli.Execute(
		t.Context(),
		[]string{"del-all", "v1.2.3"},
		directory,
		&output,
	)

	var remoteFailure *git.RemoteError
	if !errors.As(err, &remoteFailure) || remoteFailure.Remote != "a-rejected" {
		t.Fatalf("error=%v, want rejected remote", err)
	}

	var gitFailure *git.GitCommandError
	if !errors.As(err, &gitFailure) || len(gitFailure.ArgumentList) == 0 || gitFailure.ArgumentList[0] != "push" {
		t.Fatalf("error=%v, want underlying Git push error", err)
	}

	if !errors.Is(err, gitFailure.Cause) {
		t.Fatalf("error=%v, want original Git failure in error chain", err)
	}

	command.AssertTagPresence(
		t,
		directory,
		"v1.2.3",
		true,
	)
	command.AssertTagPresence(
		t,
		rejected,
		"v1.2.3",
		true,
	)
	command.AssertTagPresence(
		t,
		working,
		"v1.2.3",
		false,
	)
}
