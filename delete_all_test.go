package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDeleteAllTags(t *testing.T) {
	t.Parallel()
	directory := createRepository(t)
	runGitCommand(t, directory, "tag", "v1.2.3")
	runGitCommand(t, directory, "tag", "v1.2.4")

	remotes := []string{createBareRemote(t, directory, "origin"), createBareRemote(t, directory, "upstream")}
	for _, remote := range remotes {
		runGitCommand(t, directory, "push", remote, "refs/tags/v1.2.3", "refs/tags/v1.2.4")
	}

	var output bytes.Buffer

	err := runApplication(t.Context(), []string{"del-all", "v1.2.3"}, directory, &output)
	if err != nil {
		t.Fatal(err)
	}

	for _, repository := range append(remotes, directory) {
		assertTagPresence(t, repository, "v1.2.3", false)
		assertTagPresence(t, repository, "v1.2.4", true)
	}

	for _, name := range []string{"origin:", "upstream:", "local:"} {
		if !strings.Contains(output.String(), name) {
			t.Errorf("output = %q, missing %q", output.String(), name)
		}
	}

	err = runApplication(t.Context(), []string{"del-all", "v1.2.3"}, directory, &output)
	if err != nil {
		t.Fatalf("retry: %v", err)
	}
}

func TestDeleteAllTagsPartialFailureAndRetry(t *testing.T) {
	t.Parallel()
	directory := createRepository(t)
	runGitCommand(t, directory, "tag", "v1.2.3")
	missing := filepath.Join(t.TempDir(), "missing.git")
	runGitCommand(t, directory, "remote", "add", "a-broken", missing)
	remote := createBareRemote(t, directory, "z-working")
	runGitCommand(t, directory, "push", remote, "refs/tags/v1.2.3")

	var stdout, stderr bytes.Buffer

	code := runCommand(t.Context(), []string{"del-all", "v1.2.3"}, directory, &stdout, &stderr)
	if code != 1 || !strings.Contains(stderr.String(), "a-broken") || !strings.Contains(stdout.String(), "z-working:") {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}

	assertTagPresence(t, directory, "v1.2.3", true)
	assertTagPresence(t, remote, "v1.2.3", false)
	runGitCommand(t, directory, "init", "--bare", missing)
	runGitCommand(t, directory, "push", missing, "refs/tags/v1.2.3")

	err := runApplication(t.Context(), []string{"del-all", "v1.2.3"}, directory, &stdout)
	if err != nil {
		t.Fatal(err)
	}

	assertTagPresence(t, directory, "v1.2.3", false)
	assertTagPresence(t, missing, "v1.2.3", false)
}

func TestDeleteAllTagsWithoutRemotes(t *testing.T) {
	t.Parallel()
	directory := createRepository(t)
	runGitCommand(t, directory, "tag", "v1.2.3")

	var output bytes.Buffer
	for range 2 {
		err := runApplication(t.Context(), []string{"del-all", "v1.2.3"}, directory, &output)
		if err != nil {
			t.Fatal(err)
		}

		assertTagPresence(t, directory, "v1.2.3", false)
	}
}

func TestDeleteAllTagsRejectsInvalidArguments(t *testing.T) {
	t.Parallel()
	directory := createRepository(t)
	runGitCommand(t, directory, "tag", "v1.2.3")

	for _, args := range [][]string{{"del-all"}, {"del-all", "1.2.3"}, {"del-all", "v01.2.3"}, {"del-all", "v1.2.3", "extra"}} {
		var stdout, stderr bytes.Buffer
		if code := runCommand(t.Context(), args, directory, &stdout, &stderr); code != 1 || stdout.Len() != 0 {
			t.Fatalf("args=%q code=%d stdout=%q", args, code, stdout.String())
		}

		assertTagPresence(t, directory, "v1.2.3", true)
	}
}

func TestDeleteAllTagsUsesAllPushURLs(t *testing.T) {
	t.Parallel()
	directory := createRepository(t)
	runGitCommand(t, directory, "tag", "v1.2.3")
	fetch := createBareRemote(t, directory, "origin")
	first := t.TempDir()

	second := t.TempDir()
	for _, remote := range []string{fetch, first, second} {
		runGitCommand(t, directory, "init", "--bare", remote)
		runGitCommand(t, directory, "push", remote, "refs/tags/v1.2.3")
	}

	for _, remote := range []string{first, second} {
		runGitCommand(t, directory, "remote", "set-url", "--add", "--push", "origin", remote)
	}

	var output bytes.Buffer

	err := runApplication(t.Context(), []string{"del-all", "v1.2.3"}, directory, &output)
	if err != nil {
		t.Fatal(err)
	}

	for _, repository := range []string{first, second, directory} {
		assertTagPresence(t, repository, "v1.2.3", false)
	}

	assertTagPresence(t, fetch, "v1.2.3", true)
}

func createBareRemote(t *testing.T, directory, name string) string {
	t.Helper()
	remote := t.TempDir()
	runGitCommand(t, directory, "init", "--bare", remote)
	runGitCommand(t, directory, "remote", "add", name, remote)

	return remote
}

func assertTagPresence(t *testing.T, directory, tag string, want bool) {
	t.Helper()

	found, err := newGit(directory).hasTag(t.Context(), tag)
	if err != nil || found != want {
		t.Fatalf("tag %s in %s: found=%t err=%v want=%t", tag, directory, found, err, want)
	}
}

func TestDeleteAllTagsHandlesRemoteOnlyTag(t *testing.T) {
	t.Parallel()
	directory := createRepository(t)
	runGitCommand(t, directory, "tag", "v1.2.3")
	remote := createBareRemote(t, directory, "origin")
	runGitCommand(t, directory, "push", remote, "refs/tags/v1.2.3")
	runGitCommand(t, directory, "tag", "--delete", "v1.2.3")

	var output bytes.Buffer

	err := runApplication(t.Context(), []string{"del-all", "v1.2.3"}, directory, &output)
	if err != nil {
		t.Fatal(err)
	}

	assertTagPresence(t, remote, "v1.2.3", false)
}

func TestDeleteAllTagsHandlesRejectedPush(t *testing.T) {
	t.Parallel()
	directory := createRepository(t)
	runGitCommand(t, directory, "tag", "v1.2.3")
	rejected := createBareRemote(t, directory, "a-rejected")

	working := createBareRemote(t, directory, "z-working")
	for _, remote := range []string{rejected, working} {
		runGitCommand(t, directory, "push", remote, "refs/tags/v1.2.3")
	}

	err := os.WriteFile(filepath.Join(rejected, "hooks", "pre-receive"), []byte("#!/bin/sh\nexit 1\n"), 0o700)
	if err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer

	err = runApplication(t.Context(), []string{"del-all", "v1.2.3"}, directory, &output)

	var remoteFailure *remoteError
	if !errors.As(err, &remoteFailure) || remoteFailure.remote != "a-rejected" {
		t.Fatalf("error=%v, want rejected remote", err)
	}

	var gitFailure *gitCommandError
	if !errors.As(err, &gitFailure) || len(gitFailure.argumentList) == 0 || gitFailure.argumentList[0] != "push" {
		t.Fatalf("error=%v, want underlying Git push error", err)
	}

	if !errors.Is(err, gitFailure.cause) {
		t.Fatalf("error=%v, want original Git failure in error chain", err)
	}

	assertTagPresence(t, directory, "v1.2.3", true)
	assertTagPresence(t, rejected, "v1.2.3", true)
	assertTagPresence(t, working, "v1.2.3", false)
}
