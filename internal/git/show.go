package git

import (
	"context"
	"strings"
)

type TagInfo struct {
	Commit     string
	CommitDate string
	Message    string
}

func (g Git) GetTagInfo(ctx context.Context, tag string) (info TagInfo, err error) {
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

	info.CommitDate, err = g.RunCommand(
		ctx,
		"show",
		"--no-patch",
		"--no-notes",
		"--no-show-signature",
		"--format=%cI",
		info.Commit,
		"--",
	)
	if err != nil {
		return info, err
	}

	info.CommitDate = strings.TrimSpace(info.CommitDate)

	var objectType string

	objectType, err = g.RunCommand(
		ctx,
		"cat-file",
		"-t",
		ref,
	)
	if err != nil {
		return info, err
	}

	if strings.TrimSpace(objectType) != "tag" {
		return info, nil
	}

	var object string

	object, err = g.RunCommand(
		ctx,
		"cat-file",
		"-p",
		ref,
	)
	if err != nil {
		return info, err
	}

	_, info.Message, _ = strings.Cut(object, "\n\n")
	info.Message = strings.TrimRight(info.Message, "\n")

	return info, nil
}
