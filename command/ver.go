package command

import (
	"fmt"
	"io"
)

// Set by the linker using -X when building a release.
var releaseVersion = "dev" //nolint:gochecknoglobals // Linker injection requires a package variable.

func Ver(argumentList []string, stdout io.Writer) (err error) {
	if len(argumentList) != 0 {
		return NewUsageError()
	}

	_, err = fmt.Fprintln(stdout, releaseVersion)

	return err
}
