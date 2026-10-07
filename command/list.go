package command

import (
	"context"
	"io"
	"strconv"

	"github.com/leprosus/tagit/internal/git"
	"github.com/leprosus/tagit/internal/version"
)

func List(
	ctx context.Context,
	curGit git.Git,
	argumentList []string,
	stdout io.Writer,
) (err error) {
	var limit int

	limit, err = parseListLimit(argumentList)
	if err != nil {
		return err
	}

	var versionList version.VersionList

	versionList, err = curGit.GetSortedVersionList(ctx)
	if err != nil {
		return err
	}

	return versionList.Last(limit).Write(stdout)
}

func parseListLimit(argumentList []string) (limit int, err error) {
	if len(argumentList) > 1 {
		return limit, NewUsageError()
	}

	if len(argumentList) == 0 {
		return limit, nil
	}

	limit, err = strconv.Atoi(argumentList[0])
	if err != nil || limit < 1 {
		return limit, NewInvalidListLimitError(argumentList[0])
	}

	return limit, nil
}
