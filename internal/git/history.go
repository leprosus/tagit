package git

import (
	"context"
	"slices"
	"strings"
)

func (g Git) GetHistory(ctx context.Context, tag string) (commits CommitInfoList, err error) {
	commits, err = g.getHistory(ctx, tag)
	if err != nil || len(commits) == 0 {
		return commits, err
	}

	var tags map[string][]string

	tags, err = g.commitTags(ctx)
	if err != nil {
		return nil, err
	}

	slices.Reverse(commits)

	for i := range commits {
		commits[i].Tag = "(none)"
		if names := tags[commits[i].Commit]; len(names) > 0 {
			commits[i].Tag = strings.Join(names, ", ")
		}
	}

	return commits, nil
}

func (g Git) getHistory(ctx context.Context, tag string) (commits CommitInfoList, err error) {
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

func (g Git) commitTags(ctx context.Context) (tags map[string][]string, err error) {
	var output string

	output, err = g.RunCommand(
		ctx,
		"for-each-ref",
		"--sort=refname",
		"--format=%(objectname) %(*objectname) %(refname:strip=2)",
		"refs/tags/",
	)
	if err != nil {
		return nil, err
	}

	const (
		lightweightFields = 2
		annotatedFields   = 3
	)

	tags = make(map[string][]string)

	for line := range strings.SplitSeq(strings.TrimSuffix(output, "\n"), "\n") {
		fields := strings.Fields(line)
		if len(fields) < lightweightFields {
			continue
		}

		hash := fields[0]
		if len(fields) == annotatedFields {
			hash = fields[1]
		}

		name := fields[len(fields)-1]
		tags[hash] = append(tags[hash], name)
	}

	return tags, nil
}
