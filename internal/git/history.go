package git

import (
	"context"
	"strings"
)

func (g Git) GetHistory(ctx context.Context, tag string) (commits CommitInfoList, err error) {
	var revision string

	revision, err = g.historyRevision(ctx, tag)
	if err != nil {
		return nil, err
	}

	var output string

	output, err = g.RunCommand(
		ctx,
		"log",
		"--no-color",
		"--no-decorate",
		"--no-notes",
		"--no-show-signature",
		"--date-order",
		"-z",
		"--format=%H%x00%an <%ae>%x00%cI%x00%B",
		revision,
		"--",
	)
	if err != nil {
		return nil, err
	}

	if output == "" {
		return nil, nil
	}

	const fieldCount = 4

	// NUL separates fields and records; newlines belong to the commit message.
	fields := strings.Split(strings.TrimSuffix(output, "\x00"), "\x00")
	if len(fields)%fieldCount != 0 {
		return nil, NewInvalidHistoryOutputError()
	}

	commits = make(CommitInfoList, 0, len(fields)/fieldCount)
	for i := 0; i < len(fields); i += fieldCount {
		commits = append(commits, CommitInfo{
			Commit:  fields[i],
			Author:  fields[i+1],
			Date:    fields[i+2],
			Message: fields[i+3],
		})
	}

	return commits, nil
}
