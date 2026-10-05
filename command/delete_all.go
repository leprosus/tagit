package command

import (
	"context"
	"io"

	errtypes "github.com/leprosus/tagit/internal/errors"
	"github.com/leprosus/tagit/internal/git"
	"github.com/leprosus/tagit/internal/version"
)

func DeleteAll(ctx context.Context, curGit git.Git, argumentList []string, stdout io.Writer) (err error) {
	if len(argumentList) != 1 {
		return errtypes.NewUsageError()
	}

	if _, isValid := version.Parse(argumentList[0]); !isValid {
		return errtypes.NewUsageError()
	}

	return curGit.DeleteAllTags(ctx, argumentList[0], stdout)
}
