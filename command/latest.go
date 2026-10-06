package command

import (
	"context"
	"fmt"
	"io"

	"github.com/leprosus/tagit/internal/git"
	"github.com/leprosus/tagit/internal/version"
)

func Latest(
	ctx context.Context,
	curGit git.Git,
	argumentList []string,
	stdout io.Writer,
) (err error) {
	if len(argumentList) != 0 {
		return NewUsageError()
	}

	var (
		latest version.Version
		found  bool
	)

	latest, found, err = curGit.GetLatestVersion(ctx)
	if err != nil {
		return err
	}

	if !found {
		return NewNoVersionTagsError()
	}

	_, err = fmt.Fprintln(stdout, latest)

	return err
}
