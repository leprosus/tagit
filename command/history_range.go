package command

import (
	"context"

	"github.com/leprosus/tagit/internal/git"
	"github.com/leprosus/tagit/internal/version"
)

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
