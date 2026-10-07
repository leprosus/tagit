package command

import (
	"context"
	"io"

	"github.com/leprosus/tagit/internal/git"
	"github.com/leprosus/tagit/internal/version"
)

func Changes(
	ctx context.Context,
	curGit git.Git,
	argumentList []string,
	stdout io.Writer,
) (err error) {
	if len(argumentList) > 1 {
		return NewUsageError()
	}

	tag := ""

	if len(argumentList) == 1 {
		_, valid := version.Parse(argumentList[0])
		if !valid {
			return NewUsageError()
		}

		tag = argumentList[0]
	} else {
		var (
			latest version.Version
			found  bool
		)

		latest, found, err = curGit.GetLatestVersion(ctx)
		if err != nil {
			return err
		}

		if found {
			tag = latest.String()
		}
	}

	var output string

	output, err = curGit.GetChanges(ctx, tag)
	if err != nil {
		return err
	}

	_, err = io.WriteString(stdout, output)

	return err
}
