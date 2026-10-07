package git

import "context"

// GetChanges returns commits reachable from HEAD but not from tag.
// An empty tag selects the entire history reachable from HEAD.
func (g Git) GetChanges(ctx context.Context, tag string) (output string, err error) {
	revision := "HEAD"

	if tag != "" {
		var found bool

		found, err = g.HasTag(ctx, tag)
		if err != nil {
			return "", err
		}

		if !found {
			return "", NewTagNotFoundError(tag)
		}

		revision = "refs/tags/" + tag + "..HEAD"
	}

	return g.RunCommand(
		ctx,
		"log",
		"--no-color",
		"--no-decorate",
		"--no-notes",
		"--no-show-signature",
		"--date-order",
		"--format=%h %s",
		revision,
		"--",
	)
}
