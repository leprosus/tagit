package command

import (
	"context"
	"fmt"
	"io"

	errtypes "github.com/leprosus/tagit/internal/errors"
	"github.com/leprosus/tagit/internal/git"
	"github.com/leprosus/tagit/internal/version"
)

func Set(
	ctx context.Context,
	curGit git.Git,
	argumentList []string,
	stdout io.Writer,
) (err error) {
	if len(argumentList) != 1 {
		return errtypes.NewUsageError()
	}

	_, isValid := version.Parse(argumentList[0])
	if !isValid {
		return errtypes.NewUsageError()
	}

	err = curGit.CreateTag(ctx, argumentList[0])
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(stdout, argumentList[0])
	if err != nil {
		return err
	}

	return nil
}
