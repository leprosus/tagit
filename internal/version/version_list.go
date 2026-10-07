package version

import (
	"fmt"
	"io"
	"slices"
)

type VersionList []Version

func ParseList(tags []string) (versions VersionList) {
	for _, tag := range tags {
		current, valid := Parse(tag)
		if valid {
			versions = append(versions, current)
		}
	}

	return versions
}

func (vl VersionList) Write(stdout io.Writer) (err error) {
	for _, current := range vl {
		_, err = fmt.Fprintln(stdout, current)
		if err != nil {
			return err
		}
	}

	return nil
}

// Last returns a view of the last limit entries in an ascending sorted list.
// A non-positive limit returns the entire list.
func (vl VersionList) Last(limit int) (result VersionList) {
	if limit <= 0 || len(vl) <= limit {
		return vl
	}

	return vl[len(vl)-limit:]
}

// Sort orders the list in place from smallest to greatest version.
func (vl VersionList) Sort() {
	slices.SortFunc(vl, func(left, right Version) int { return left.Compare(right) })
}

// Latest finds the greatest version without sorting or modifying the list.
func (vl VersionList) Latest() (latest Version, found bool) {
	for _, current := range vl {
		if !found || current.Compare(latest) > 0 {
			latest = current
			found = true
		}
	}

	return latest, found
}
