package git

import (
	"context"
	"sort"

	"github.com/leprosus/tagit/internal/version"
)

func (g Git) GetSortedVersionList(ctx context.Context) (versionList []version.Version, err error) {
	var tagList []string

	tagList, err = g.GetTagList(ctx)
	if err != nil {
		return versionList, err
	}

	for _, tag := range tagList {
		currentVersion, isValid := version.Parse(tag)
		if !isValid {
			continue
		}

		versionList = append(versionList, currentVersion)
	}

	sort.Slice(versionList, func(left, right int) bool {
		return versionList[left].Compare(versionList[right]) < 0
	})

	return versionList, nil
}
