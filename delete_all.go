package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
)

func deleteAllVersionTags(ctx context.Context, curGit git, argumentList []string, stdout io.Writer) (err error) {
	if len(argumentList) != 1 {
		return newUsageError()
	}

	if _, isValid := parseVersion(argumentList[0]); !isValid {
		return newUsageError()
	}

	return curGit.deleteAllTags(ctx, argumentList[0], stdout)
}

func (g git) deleteAllTags(ctx context.Context, tag string, stdout io.Writer) (err error) {
	var output string

	output, err = g.runCommand(ctx, "remote")
	if err != nil {
		return err
	}

	var failures []error

	for remote := range strings.FieldsSeq(output) {
		err = g.deleteRemoteTag(ctx, remote, tag)
		if err != nil {
			failures = append(failures, newRemoteError(remote, err))

			continue
		}

		_, err = fmt.Fprintf(stdout, "%s: %s deleted or already absent\n", remote, tag)
		if err != nil {
			failures = append(failures, err)
		}
	}

	if len(failures) != 0 {
		return errors.Join(failures...)
	}

	err = g.deleteTag(ctx, tag)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(stdout, "local: %s deleted or already absent\n", tag)
	if err != nil {
		return err
	}

	return nil
}

func (g git) deleteRemoteTag(ctx context.Context, remote, tag string) (err error) {
	urls, err := g.runCommand(ctx, "remote", "get-url", "--push", "--all", remote)
	if err != nil {
		return err
	}

	var failures []error

	for url := range strings.SplitSeq(strings.TrimSpace(urls), "\n") {
		err = g.deleteTagAtURL(ctx, url, tag)
		if err != nil {
			failures = append(failures, err)
		}
	}

	return errors.Join(failures...)
}

func (g git) deleteTagAtURL(ctx context.Context, url, tag string) (err error) {
	ref := "refs/tags/" + tag

	var output string

	output, err = g.runCommand(ctx, "ls-remote", "--refs", "--tags", "--", url, ref)
	if err != nil {
		return err
	}

	for line := range strings.SplitSeq(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 || fields[1] != ref {
			continue
		}

		_, err = g.runCommand(ctx, "push", "--no-follow-tags", "--", url, ":"+ref)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}
