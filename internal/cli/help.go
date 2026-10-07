package cli

import (
	"fmt"
	"io"

	"github.com/leprosus/tagit/command"
)

const HelpMessage = `Usage: tagit [patch|minor|major]
       tagit help [command]
       tagit ver
       tagit changes [vMAJOR.MINOR.PATCH]
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
	if name == "patch" || name == "minor" || name == "major" {
		return "Usage: tagit " + name + "\n\nCreates the next local version tag using the " + name + " increment.\n"
	}

	messages := map[string]string{
		"changes": "Usage: tagit changes [vMAJOR.MINOR.PATCH]\n\n" +
			"Lists commits after the specified local tag up to HEAD, newest first.\n" +
			"Defaults to the greatest local version tag, or all HEAD history if none exist.\n",
		"show": "Usage: tagit show vMAJOR.MINOR.PATCH\n\n" +
			"Shows the local tag's commit, commit date, and annotated tag message when present.\n",
		"latest": "Usage: tagit latest\n\nPrints the greatest local version tag. Exits with status 1 if none exist.\n",
		"list": "Usage: tagit list [limit]\n\n" +
			"Lists version tags from oldest to newest. The optional limit must be a positive integer.\n",
		"set":     "Usage: tagit set vMAJOR.MINOR.PATCH\n\nCreates the specified local version tag.\n",
		"del":     "Usage: tagit del vMAJOR.MINOR.PATCH\n\nDeletes the local tag when it exists.\n",
		"del-all": "Usage: tagit del-all vMAJOR.MINOR.PATCH\n\nDeletes the tag from all remote push destinations, then locally.\n",
		"ver": "Usage: tagit ver\n\n" +
			"Prints the executable's release version, or dev for a local build. Does not require Git.\n",
	}

	return messages[name]
}
