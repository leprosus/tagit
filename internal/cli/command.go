package cli

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/leprosus/tagit/command"
	"github.com/leprosus/tagit/internal/git"
	"github.com/leprosus/tagit/internal/version"
)

func Execute(
	ctx context.Context,
	argumentList []string,
	dirPath string,
	stdout io.Writer,
) (err error) {
	if len(argumentList) > 0 && isHelpRequest(argumentList[0]) {
		return showHelp(argumentList[1:], stdout)
	}

	if len(argumentList) == 2 && isHelpFlag(argumentList[1]) {
		return showHelp(argumentList[:1], stdout)
	}

	if len(argumentList) > 0 && argumentList[0] == "ver" {
		return command.Ver(argumentList[1:], stdout)
	}

	curGit := git.New(dirPath)

	var isRepository bool

	isRepository, err = curGit.IsRepository(ctx)
	if err != nil {
		return err
	}

	if !isRepository {
		return printHelp(stdout)
	}

	if len(argumentList) == 0 {
		return command.Increase(
			ctx,
			curGit,
			version.Patch,
			stdout,
		)
	}

	return dispatchCommand(
		ctx,
		curGit,
		argumentList,
		stdout,
	)
}

func dispatchCommand(
	ctx context.Context,
	curGit git.Git,
	argumentList []string,
	stdout io.Writer,
) (err error) {
	commandName, arg := argumentList[0], argumentList[1:]

	handlers := map[string]func(
		context.Context,
		git.Git,
		[]string,
		io.Writer,
	) error{
		"history": command.History,
		"changes": command.Changes,
		"show":    command.Show,
		"latest":  command.Latest,
		"list":    command.List,
	}

	handler, found := handlers[commandName]
	if found {
		return handler(
			ctx,
			curGit,
			arg,
			stdout,
		)
	}

	switch commandName {
	case "set":
		return setTag(
			ctx,
			curGit,
			arg,
			stdout,
		)
	case "del":
		return command.Delete(ctx, curGit, arg)
	case "del-all":
		return command.DeleteAll(
			ctx,
			curGit,
			arg,
			stdout,
		)
	}

	if len(argumentList) > 1 {
		return command.NewUsageError()
	}

	incrementKind := version.Kind(argumentList[0])

	if !incrementKind.IsValid() {
		return command.NewUnknownVersionIncrementError(string(incrementKind))
	}

	return command.Increase(
		ctx,
		curGit,
		incrementKind,
		stdout,
	)
}

// Run executes the CLI and returns its exit code.
func Run(
	ctx context.Context,
	argumentList []string,
	dirPath string,
	stdout io.Writer,
	stderr io.Writer,
) (exitCode int) {
	err := Execute(
		ctx,
		argumentList,
		dirPath,
		stdout,
	)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, "tagit:", err)
		_ = printHelp(stderr)

		return 1
	}

	return 0
}

func setTag(
	ctx context.Context,
	curGit git.Git,
	argumentList []string,
	stdout io.Writer,
) (err error) {
	err = command.Set(
		ctx,
		curGit,
		argumentList,
		stdout,
	)

	_, isUsageError := errors.AsType[*command.UsageError](err)
	if isUsageError {
		return printHelp(stdout)
	}

	return err
}
