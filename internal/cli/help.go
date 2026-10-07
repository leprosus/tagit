package cli

import (
	"fmt"
	"io"

	"github.com/leprosus/tagit/command"
)

const HelpMessage = `Usage: tagit [patch|minor|major]
       tagit help [command]
       tagit ver
       tagit latest
       tagit show vMAJOR.MINOR.PATCH
       tagit list [limit]
       tagit set vMAJOR.MINOR.PATCH
       tagit del vMAJOR.MINOR.PATCH
       tagit del-all vMAJOR.MINOR.PATCH

Creates, increments, lists, or deletes semantic Git tags.
del-all deletes a tag from all remotes before deleting it locally.
Use -h or --help to show help. Commands also accept -h or --help.
`

func printHelp(stdout io.Writer) (err error) {
	_, err = fmt.Fprint(stdout, HelpMessage)

	return err
}

func isHelpFlag(argument string) (result bool) {
	return argument == "-h" || argument == "--help"
}

func isHelpRequest(argument string) (result bool) {
	return argument == "help" || isHelpFlag(argument)
}

func showHelp(argumentList []string, stdout io.Writer) (err error) {
	if len(argumentList) == 0 {
		return printHelp(stdout)
	}

	if len(argumentList) != 1 {
		return command.NewUsageError()
	}

	message := commandHelp(argumentList[0])
	if message == "" {
		return command.NewUsageError()
	}

	_, err = fmt.Fprint(stdout, message)

	return err
}

func commandHelp(name string) (message string) {
	switch name {
	case "show":
		return "Usage: tagit show vMAJOR.MINOR.PATCH\n\n" +
			"Shows the local tag's commit, commit date, and annotated tag message when present.\n"
	case "latest":
		return "Usage: tagit latest\n\nPrints the greatest local version tag. Exits with status 1 if none exist.\n"
	case "patch", "minor", "major":
		return "Usage: tagit " + name + "\n\nCreates the next local version tag using the " + name + " increment.\n"
	case "list":
		return "Usage: tagit list [limit]\n\n" +
			"Lists version tags from oldest to newest. The optional limit must be a positive integer.\n"
	case "set":
		return "Usage: tagit set vMAJOR.MINOR.PATCH\n\nCreates the specified local version tag.\n"
	case "del":
		return "Usage: tagit del vMAJOR.MINOR.PATCH\n\nDeletes the local tag when it exists.\n"
	case "del-all":
		return "Usage: tagit del-all vMAJOR.MINOR.PATCH\n\nDeletes the tag from all remote push destinations, then locally.\n"
	case "ver":
		return "Usage: tagit ver\n\n" +
			"Prints the executable's release version, or dev for a local build. Does not require Git.\n"
	default:
		return ""
	}
}
