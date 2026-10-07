package command

import (
	"context"
	"io"

	"github.com/leprosus/tagit/internal/git"
	"github.com/leprosus/tagit/internal/version"
)

func History(
	ctx context.Context,
	curGit git.Git,
	argumentList []string,
	stdout io.Writer,
) (err error) {
	var tag string

	tag, err = historyTag(ctx, curGit, argumentList)
	if err != nil {
		return err
	}

	var commits git.CommitInfoList

	commits, err = curGit.GetHistory(ctx, tag)
	if err != nil {
		return err
	}

	return commits.Write(stdout)
}

func historyTag(ctx context.Context, curGit git.Git, argumentList []string) (tag string, err error) {
	if len(argumentList) > 1 {
		return "", NewUsageError()
	}

	if len(argumentList) == 1 {
		_, valid := version.Parse(argumentList[0])
		if !valid {
			return "", NewUsageError()
		}

		return argumentList[0], nil
	}

	var (
		latest version.Version
		found  bool
	)

	latest, found, err = curGit.GetLatestVersion(ctx)
	if err != nil {
		return "", err
	}

	if found {
		return latest.String(), nil
	}

	return "", nil
}
