package main

import (
	"context"
	"fmt"
	"io"
)

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

func (g git) increaseTag(ctx context.Context, kind kind) (tag string, err error) {
	var tagList []string

	tagList, err = g.getTagList(ctx)
	if err != nil {
		return tag, err
	}

	current, found := getLatestVersion(tagList)
	if !found {
		current = version{}

		kind = patch
	}

	var next version

	next, err = current.getNextVersion(kind)
	if err != nil {
		return tag, err
	}

	tag = next.String()

	err = g.createTag(ctx, tag)
	if err != nil {
		return tag, err
	}

	return tag, nil
}
