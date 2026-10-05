package main

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strconv"
)

func printVersionTagList(ctx context.Context, curGit git, argumentList []string, stdout io.Writer) (err error) {
	var limit int

	limit, err = parseListLimit(argumentList)
	if err != nil {
		return err
	}

	var versionList []version

	versionList, err = curGit.getSortedVersionList(ctx)
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
		return limit, newUsageError()
	}

	if len(argumentList) == 0 {
		return limit, nil
	}

	limit, err = strconv.Atoi(argumentList[0])
	if err != nil || limit < 1 {
		return limit, newInvalidListLimitError(argumentList[0])
	}

	return limit, nil
}

func truncateVersionList(versionList []version, limit int) (result []version) {
	if limit == 0 || len(versionList) <= limit {
		return versionList
	}

	result = versionList[len(versionList)-limit:]

	return result
}

func printVersionList(versionList []version, stdout io.Writer) (err error) {
	for _, currentVersion := range versionList {
		_, err = fmt.Fprintln(stdout, currentVersion)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g git) getSortedVersionList(ctx context.Context) (versionList []version, err error) {
	var tagList []string

	tagList, err = g.getTagList(ctx)
	if err != nil {
		return versionList, err
	}

	for _, tag := range tagList {
		currentVersion, isValid := parseVersion(tag)
		if !isValid {
			continue
		}

		versionList = append(versionList, currentVersion)
	}

	sort.Slice(versionList, func(left, right int) bool {
		return versionList[left].Compare(versionList[right]) < 0
	})

	return versionList, nil
}
