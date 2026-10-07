package git

import (
	"context"

	"github.com/leprosus/tagit/internal/version"
)

func (g Git) GetLatestVersion(ctx context.Context) (latest version.Version, found bool, err error) {
	var tagList []string

	tagList, err = g.GetTagList(ctx)
	if err != nil {
		return latest, false, err
	}

	latest, found = version.ParseList(tagList).Latest()

	return latest, found, nil
}
