package main

import (
	"context"
	"os"
	"os/exec"
	"slices"
	"sort"
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

func (g git) deleteTag(ctx context.Context, tag string) (err error) {
	var hasTag bool

	hasTag, err = g.hasTag(ctx, tag)
	if err != nil || !hasTag {
		return err
	}

	_, err = g.runCommand(ctx, "tag", "--delete", tag)

	return err
}

func (g git) hasTag(ctx context.Context, tag string) (result bool, err error) {
	var tagList []string

	tagList, err = g.getTagList(ctx)
	if err != nil {
		return result, err
	}

	if slices.Contains(tagList, tag) {
		result = true

		return result, nil
	}

	return result, nil
}

func (g git) isRepository(ctx context.Context) (result bool, err error) {
	argumentList := []string{"rev-parse", "--git-dir"}
	command := makeCommand(ctx, argumentList...)
	command.Dir = g.dirPath

	const localisation = "LC_ALL=C"

	command.Env = append(os.Environ(), localisation)

	var output []byte

	output, err = command.CombinedOutput()
	if err == nil {
		return true, nil
	}

	response := strings.TrimSpace(string(output))

	const notRepositoryErrorPrefix = "fatal: not a git repository"

	isNotRepository := strings.HasPrefix(response, notRepositoryErrorPrefix)

	if command.ProcessState != nil && command.ProcessState.ExitCode() == 128 && isNotRepository {
		return false, nil
	}

	return false, newGitCommandError(argumentList, err, response)
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

	return versionList, err
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
