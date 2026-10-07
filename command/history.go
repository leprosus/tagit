package command

import (
	"context"
	"io"

	"github.com/leprosus/tagit/internal/git"
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

	return commits.WriteCommitList(stdout)
}
