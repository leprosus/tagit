package git

import (
	"context"
	"strings"
)

func (g Git) GetTagInfo(ctx context.Context, tag string) (info CommitInfo, err error) {
	var found bool

	found, err = g.HasTag(ctx, tag)
	if err != nil {
		return info, err
	}

	if !found {
		return info, NewTagNotFoundError(tag)
	}

	ref := "refs/tags/" + tag

	info.Commit, err = g.RunCommand(
		ctx,
		"rev-parse",
		"--verify",
		ref+"^{commit}",
	)
	if err != nil {
		return info, err
	}

	info.Commit = strings.TrimSpace(info.Commit)

	var metadata string

	metadata, err = g.RunCommand(
		ctx,
		"show",
		"--no-patch",
		"--no-notes",
		"--no-show-signature",
		"--format=%an <%ae>%x00%cI%x00%B",
		info.Commit,
		"--",
	)
	if err != nil {
		return info, err
	}

	info.Author, metadata, _ = strings.Cut(metadata, "\x00")
	info.Date, info.Message, _ = strings.Cut(metadata, "\x00")

	return info, nil
}
