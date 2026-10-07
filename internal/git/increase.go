package git

import (
	"context"

	"github.com/leprosus/tagit/internal/version"
)

func (g Git) IncreaseTag(ctx context.Context, incrementKind version.Kind) (tag string, err error) {
	var tagList []string

	tagList, err = g.GetTagList(ctx)
	if err != nil {
		return tag, err
	}

	current, found := version.ParseList(tagList).Latest()
	if !found {
		current = version.Version{}

		incrementKind = version.Patch
	}

	var next version.Version

	next, err = current.Next(incrementKind)
	if err != nil {
		return tag, err
	}

	tag = next.String()

	err = g.CreateTag(ctx, tag)
	if err != nil {
		return tag, err
	}

	return tag, nil
}
