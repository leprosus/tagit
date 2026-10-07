package git

import "context"

func (g Git) historyRevision(ctx context.Context, tag string) (revision string, err error) {
	revision = "HEAD"

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

	return revision, nil
}
