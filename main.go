package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
)

const helpMessage = `Usage: tagit [patch|minor|major]
       tagit list [limit]

Creates, increments, or lists semantic Git tags in the current repository.
`

func printHelp(stdout io.Writer) (err error) {
	_, err = fmt.Fprint(stdout, helpMessage)

	return err
}

func runApplication(ctx context.Context, argumentList []string, dirPath string, stdout io.Writer) (err error) {
	curGit := newGit(dirPath)

	if len(argumentList) == 0 {
		if !curGit.isRepository(ctx) {
			return printHelp(stdout)
		}

		return increaseTag(ctx, curGit, patch, stdout)
	}

	if argumentList[0] == "list" {
		return printVersionTagList(ctx, curGit, argumentList[1:], stdout)
	}

	if len(argumentList) > 1 {
		return newUsageError()
	}

	incrementKind := kind(argumentList[0])

	if !incrementKind.isValid() {
		return newUnknownVersionIncrementError(incrementKind)
	}

	return increaseTag(ctx, curGit, incrementKind, stdout)
}

func increaseTag(ctx context.Context, curGit git, incrementKind kind, stdout io.Writer) (err error) {
	var tag string

	tag, err = curGit.increaseTag(ctx, incrementKind)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(stdout, tag)
	if err != nil {
		return err
	}

	return err
}

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

	return err
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

	return err
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := runApplication(ctx, os.Args[1:], ".", os.Stdout)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "tagit:", err)

		return
	}
}
