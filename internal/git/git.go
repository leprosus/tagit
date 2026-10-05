package git

import (
	"context"
	"os"
	"os/exec"
	"slices"
	"strings"

	errtypes "github.com/leprosus/tagit/internal/errors"
)

type Git struct {
	dirPath string
}

func New(dirPath string) Git {
	return Git{dirPath: dirPath}
}

func MakeCommand(ctx context.Context, argumentList ...string) (command *exec.Cmd) {
	command = exec.CommandContext(ctx, "git")
	command.Args = append(command.Args, argumentList...)

	return command
}

func (g Git) GetTagList(ctx context.Context) (tagList []string, err error) {
	var output string

	output, err = g.RunCommand(ctx, "tag", "--list")
	if err != nil {
		return tagList, err
	}

	if strings.TrimSpace(output) == "" {
		return tagList, nil
	}

	tagList = strings.Fields(output)

	return tagList, nil
}

func (g Git) HasTag(ctx context.Context, tag string) (result bool, err error) {
	var tagList []string

	tagList, err = g.GetTagList(ctx)
	if err != nil {
		return result, err
	}

	if slices.Contains(tagList, tag) {
		result = true

		return result, nil
	}

	return result, nil
}

func (g Git) IsRepository(ctx context.Context) (result bool, err error) {
	argumentList := []string{"rev-parse", "--git-dir"}
	command := MakeCommand(ctx, argumentList...)
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

	return false, errtypes.NewGitCommandError(argumentList, err, response)
}

func (g Git) RunCommand(ctx context.Context, argumentList ...string) (output string, err error) {
	command := MakeCommand(ctx, argumentList...)
	command.Dir = g.dirPath

	var bs []byte

	bs, err = command.CombinedOutput()
	if err != nil {
		return output, errtypes.NewGitCommandError(argumentList, err, strings.TrimSpace(string(bs)))
	}

	output = string(bs)

	return output, nil
}
