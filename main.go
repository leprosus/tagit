package main

import (
	"context"
	"fmt"
	"io"
	"os"
)

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

	var tag string

	tag, err = newGit(dirPath).increaseTag(ctx, incrementKind)
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
