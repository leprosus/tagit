package main

import (
	"context"
	"fmt"
	"io"
)

func setVersionTag(ctx context.Context, curGit git, argumentList []string, stdout io.Writer) (err error) {
	if len(argumentList) != 1 {
		return printHelp(stdout)
	}

	_, isValid := parseVersion(argumentList[0])
	if !isValid {
		return printHelp(stdout)
	}

	err = curGit.createTag(ctx, argumentList[0])
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(stdout, argumentList[0])
	if err != nil {
		return err
	}

	return nil
}

func (g git) createTag(ctx context.Context, tag string) (err error) {
	_, err = g.runCommand(ctx, "tag", tag)
	if err != nil {
		return err
	}

	return nil
}
