package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
)

type version struct {
	major uint64
	minor uint64
	patch uint64
}

func (v version) String() (result string) {
	return fmt.Sprintf("v%d.%d.%d", v.major, v.minor, v.patch)
}

func (v version) Compare(other version) (result int) {
	if v.major < other.major {
		return -1
	}

	if v.major > other.major {
		return 1
	}

	if v.minor < other.minor {
		return -1
	}

	if v.minor > other.minor {
		return 1
	}

	if v.patch < other.patch {
		return -1
	}

	if v.patch > other.patch {
		return 1
	}

	return 0
}

var versionPattern = regexp.MustCompile(`^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$`)

func parseVersion(value string) (result version, ok bool) {
	const expectedParts = 4

	matchList := versionPattern.FindStringSubmatch(value)
	if matchList == nil || len(matchList) != expectedParts {
		return result, false
	}

	var err error

	result.major, err = strconv.ParseUint(matchList[1], 10, 64)
	if err != nil {
		return result, false
	}

	result.minor, err = strconv.ParseUint(matchList[2], 10, 64)
	if err != nil {
		return result, false
	}

	result.patch, err = strconv.ParseUint(matchList[3], 10, 64)
	if err != nil {
		return result, false
	}

	return result, true
}

func (v version) getNextVersion(kind kind) (result version, err error) {
	result = v

	switch kind {
	case patch:
		if v.patch == math.MaxUint64 {
			return result, newVersionOverflowError(patch)
		}

		result.patch++

	case minor:
		if v.minor == math.MaxUint64 {
			return result, newVersionOverflowError(minor)
		}

		result.minor++
		result.patch = 0

	case major:
		if v.major == math.MaxUint64 {
			return result, newVersionOverflowError(major)
		}

		result.major++
		result.minor = 0

	default:
		return result, newUnknownVersionIncrementError(kind)
	}

	return result, nil
}

func getLatestVersion(tagList []string) (latest version, found bool) {
	var (
		current version
		isValid bool
	)

	for _, tag := range tagList {
		current, isValid = parseVersion(tag)
		if !isValid || current.Compare(latest) <= 0 {
			continue
		}

		latest = current
		found = true
	}

	return latest, found
}
