package command

import (
	"context"
	"fmt"
	"io"

	"github.com/leprosus/tagit/internal/git"
	"github.com/leprosus/tagit/internal/version"
)

func Increase(ctx context.Context, curGit git.Git, incrementKind version.Kind, stdout io.Writer) (err error) {
	var tag string

	tag, err = curGit.IncreaseTag(ctx, incrementKind)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(stdout, tag)
	if err != nil {
		return err
	}

	return err
}
