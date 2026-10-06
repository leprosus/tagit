package git

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
)

func (g Git) DeleteAllTags(ctx context.Context, tag string, stdout io.Writer) (err error) {
	var output string

	output, err = g.RunCommand(ctx, "remote")
	if err != nil {
		return err
	}

	var failures []error

	for remote := range strings.FieldsSeq(output) {
		err = g.deleteRemoteTag(ctx, remote, tag)
		if err != nil {
			failures = append(failures, NewRemoteError(remote, err))

			continue
		}

		_, err = fmt.Fprintf(
			stdout,
			"%s: %s deleted or already absent\n",
			remote,
			tag,
		)
		if err != nil {
			failures = append(failures, err)
		}
	}

	if len(failures) != 0 {
		return errors.Join(failures...)
	}

	err = g.DeleteTag(ctx, tag)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(stdout, "local: %s deleted or already absent\n", tag)
	if err != nil {
		return err
	}

	return nil
}

func (g Git) deleteRemoteTag(ctx context.Context, remote, tag string) (err error) {
	var urls string

	urls, err = g.RunCommand(
		ctx,
		"remote",
		"get-url",
		"--push",
		"--all",
		remote,
	)
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

func (g Git) deleteTagAtURL(ctx context.Context, url, tag string) (err error) {
	ref := "refs/tags/" + tag

	var output string

	output, err = g.RunCommand(
		ctx,
		"ls-remote",
		"--refs",
		"--tags",
		"--",
		url,
		ref,
	)
	if err != nil {
		return err
	}

	for line := range strings.SplitSeq(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 || fields[1] != ref {
			continue
		}

		_, err = g.RunCommand(
			ctx,
			"push",
			"--no-follow-tags",
			"--",
			url,
			":"+ref,
		)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}
