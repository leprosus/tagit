package main

import (
	"context"
	"fmt"
	"io"
	"os"
)

const helpMessage = `Usage: tagit [patch|minor|major]
       tagit list [limit]
       tagit set vMAJOR.MINOR.PATCH
       tagit del vMAJOR.MINOR.PATCH
       tagit del-all vMAJOR.MINOR.PATCH

Creates, increments, lists, or deletes semantic Git tags.
del-all deletes a tag from all remotes before deleting it locally.
`

func printHelp(stdout io.Writer) (err error) {
	_, err = fmt.Fprint(stdout, helpMessage)

	return err
}

func runApplication(ctx context.Context, argumentList []string, dirPath string, stdout io.Writer) (err error) {
	curGit := newGit(dirPath)

	isRepository, err := curGit.isRepository(ctx)
	if err != nil {
		return err
	}

	if !isRepository {
		return printHelp(stdout)
	}

	if len(argumentList) == 0 {
		return increaseTag(ctx, curGit, patch, stdout)
	}

	command, arg := argumentList[0], argumentList[1:]

	switch command {
	case "list":
		return printVersionTagList(ctx, curGit, arg, stdout)
	case "set":
		return setVersionTag(ctx, curGit, arg, stdout)
	case "del":
		return deleteVersionTag(ctx, curGit, arg)
	case "del-all":
		return deleteAllVersionTags(ctx, curGit, arg, stdout)
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

func main() {
	exitCode := runCommand(context.Background(), os.Args[1:], ".", os.Stdout, os.Stderr)

	os.Exit(exitCode)
}

func runCommand(ctx context.Context, argumentList []string, dirPath string, stdout io.Writer, stderr io.Writer) (exitCode int) {
	err := runApplication(ctx, argumentList, dirPath, stdout)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, "tagit:", err)
		_ = printHelp(stderr)

		return 1
	}

	return 0
}
