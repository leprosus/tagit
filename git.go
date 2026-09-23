package main

import (
	"context"
	"os/exec"
	"strings"
)

type git struct {
	dirPath string
}

func newGit(dirPath string) git {
	return git{dirPath: dirPath}
}

func makeCommand(ctx context.Context, argumentList ...string) (command *exec.Cmd) {
	command = exec.CommandContext(ctx, "git")
	command.Args = append(command.Args, argumentList...)

	return command
}

func (g git) getTagList(ctx context.Context) (tagList []string, err error) {
	var output string

	output, err = g.runCommand(ctx, "tag", "--list")
	if err != nil {
		return tagList, err
	}

	if strings.TrimSpace(output) == "" {
		return tagList, nil
	}

	tagList = strings.Fields(output)

	return tagList, nil
}

func (g git) createTag(ctx context.Context, tag string) (err error) {
	_, err = g.runCommand(ctx, "tag", tag)

	return err
}

func (g git) isRepository(ctx context.Context) bool {
	var err error

	_, err = g.runCommand(ctx, "rev-parse", "--git-dir")

	return err == nil
}

func (g git) runCommand(ctx context.Context, argumentList ...string) (output string, err error) {
	command := makeCommand(ctx, argumentList...)
	command.Dir = g.dirPath

	var bs []byte

	bs, err = command.CombinedOutput()
	if err != nil {
		return output, newGitCommandError(argumentList, err, strings.TrimSpace(string(bs)))
	}

	output = string(bs)

	return output, nil
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
