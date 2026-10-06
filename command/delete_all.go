package command

import (
	"context"
	"io"

	"github.com/leprosus/tagit/internal/git"
	"github.com/leprosus/tagit/internal/version"
)

func DeleteAll(
	ctx context.Context,
	curGit git.Git,
	argumentList []string,
	stdout io.Writer,
) (err error) {
	if len(argumentList) != 1 {
		return NewUsageError()
	}

	_, isValid := version.Parse(argumentList[0])
	if !isValid {
		return NewUsageError()
	}

	return curGit.DeleteAllTags(ctx, argumentList[0], stdout)
}
