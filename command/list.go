package command

import (
	"context"
	"fmt"
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

	var versionList []version.Version

	versionList, err = curGit.GetSortedVersionList(ctx)
	if err != nil {
		return err
	}

	versionList = truncateVersionList(versionList, limit)

	err = printVersionList(versionList, stdout)
	if err != nil {
		return err
	}

	return nil
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

func truncateVersionList(versionList []version.Version, limit int) (result []version.Version) {
	if limit == 0 || len(versionList) <= limit {
		return versionList
	}

	result = versionList[len(versionList)-limit:]

	return result
}

func printVersionList(versionList []version.Version, stdout io.Writer) (err error) {
	for _, currentVersion := range versionList {
		_, err = fmt.Fprintln(stdout, currentVersion)
		if err != nil {
			return err
		}
	}

	return nil
}
