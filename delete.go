package main

import (
	"context"
)

func deleteVersionTag(ctx context.Context, curGit git, argumentList []string) (err error) {
	if len(argumentList) != 1 {
		return newUsageError()
	}

	_, isValid := parseVersion(argumentList[0])
	if !isValid {
		return newUsageError()
	}

	err = curGit.deleteTag(ctx, argumentList[0])
	if err != nil {
		return err
	}

	return nil
}

func (g git) deleteTag(ctx context.Context, tag string) (err error) {
	var hasTag bool

	hasTag, err = g.hasTag(ctx, tag)
	if err != nil || !hasTag {
		return err
	}

	_, err = g.runCommand(ctx, "tag", "--delete", tag)
	if err != nil {
		return err
	}

	return nil
}
