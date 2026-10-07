package command

import (
	"context"
	"fmt"
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

	var commits []git.CommitInfo

	commits, err = curGit.GetHistory(ctx, tag)
	if err != nil {
		return err
	}

	for i, commit := range commits {
		separator := ""
		if i > 0 {
			separator = "\n"
		}

		_, err = fmt.Fprintf(
			stdout,
			"%sCommit: %s\nAuthor: %s\nDate: %s\nMessage:\n%s\n",
			separator,
			commit.Commit,
			commit.Author,
			commit.CommitDate,
			commit.Message,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
