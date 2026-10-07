package git

import (
	"context"
	"strings"
)

type TagInfo struct {
	Commit     string
	Author     string
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

	var metadata string

	metadata, err = g.RunCommand(
		ctx,
		"show",
		"--no-patch",
		"--no-notes",
		"--no-show-signature",
		"--format=%an <%ae>%x00%cI",
		info.Commit,
		"--",
	)
	if err != nil {
		return info, err
	}

	info.Author, info.CommitDate, _ = strings.Cut(strings.TrimRight(metadata, "\n"), "\x00")

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
