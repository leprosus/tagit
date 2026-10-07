package command

import (
	"context"
	"io"

	"github.com/leprosus/tagit/internal/git"
	"github.com/leprosus/tagit/internal/version"
)

func Show(
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

	var info git.CommitInfo

	info, err = curGit.GetTagInfo(ctx, argumentList[0])
	if err != nil {
		return err
	}

	return info.Write(stdout)
}
