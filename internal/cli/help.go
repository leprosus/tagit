package cli

import (
	"fmt"
	"io"
)

const HelpMessage = `Usage: tagit [patch|minor|major]
       tagit ver
       tagit list [limit]
       tagit set vMAJOR.MINOR.PATCH
       tagit del vMAJOR.MINOR.PATCH
       tagit del-all vMAJOR.MINOR.PATCH

Creates, increments, lists, or deletes semantic Git tags.
del-all deletes a tag from all remotes before deleting it locally.
`

func printHelp(stdout io.Writer) (err error) {
	_, err = fmt.Fprint(stdout, HelpMessage)

	return err
}
