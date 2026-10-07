package command

import (
	"context"
	"fmt"
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

	var info git.TagInfo

	info, err = curGit.GetTagInfo(ctx, argumentList[0])
	if err != nil {
		return err
	}

	message := "Message: (none)"
	if info.Message != "" {
		message = "Message:\n" + info.Message
	}

	_, err = fmt.Fprintf(
		stdout,
		"Tag: %s\nCommit: %s\nDate: %s\n%s\n",
		argumentList[0],
		info.Commit,
		info.CommitDate,
		message,
	)

	return err
}
