package git

import (
	"context"

	"github.com/leprosus/tagit/internal/version"
)

func (g Git) GetSortedVersionList(ctx context.Context) (versionList version.VersionList, err error) {
	var tagList []string

	tagList, err = g.GetTagList(ctx)
	if err != nil {
		return versionList, err
	}

	versionList = version.ParseList(tagList)
	versionList.Sort()

	return versionList, nil
}
