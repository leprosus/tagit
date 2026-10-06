package command

import (
	"context"

	"github.com/leprosus/tagit/internal/git"
	"github.com/leprosus/tagit/internal/version"
)

func Delete(ctx context.Context, curGit git.Git, argumentList []string) (err error) {
	if len(argumentList) != 1 {
		return NewUsageError()
	}

	_, isValid := version.Parse(argumentList[0])
	if !isValid {
		return NewUsageError()
	}

	err = curGit.DeleteTag(ctx, argumentList[0])
	if err != nil {
		return err
	}

	return nil
}
