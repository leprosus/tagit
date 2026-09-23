package main

import (
	"context"
	"fmt"
	"io"
	"os"
)

const helpMessage = `Usage: tagit [patch|minor|major]

Creates or increments a semantic Git tag in the current repository.
`

func printHelp(stdout io.Writer) (err error) {
	_, err = fmt.Fprint(stdout, helpMessage)

	return err
}

func runApplication(argumentList []string, dirPath string, stdout io.Writer) (err error) {
	incrementKind := patch

	if len(argumentList) > 1 {
		return newUsageError()
	}

	if len(argumentList) == 1 {
		incrementKind = kind(argumentList[0])
	}

	if !incrementKind.isValid() {
		return newUnknownVersionIncrementError(incrementKind)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	curGit := newGit(dirPath)
	if len(argumentList) == 0 && !curGit.isRepository(ctx) {
		return printHelp(stdout)
	}

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

func main() {
	err := runApplication(os.Args[1:], ".", os.Stdout)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "tagit:", err)

		os.Exit(1)
	}
}
