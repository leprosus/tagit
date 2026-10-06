package version

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
)

type Version struct {
	major uint64
	minor uint64
	patch uint64
}

func (v Version) String() (result string) {
	return fmt.Sprintf(
		"v%d.%d.%d",
		v.major,
		v.minor,
		v.patch,
	)
}

func (v Version) Compare(other Version) (result int) {
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

func Parse(value string) (result Version, ok bool) {
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

func (v Version) Next(incrementKind Kind) (result Version, err error) {
	result = v

	switch incrementKind {
	case Patch:
		if v.patch == math.MaxUint64 {
			return result, NewVersionOverflowError(string(Patch))
		}

		result.patch++

	case Minor:
		if v.minor == math.MaxUint64 {
			return result, NewVersionOverflowError(string(Minor))
		}

		result.minor++
		result.patch = 0

	case Major:
		if v.major == math.MaxUint64 {
			return result, NewVersionOverflowError(string(Major))
		}

		result.major++
		result.minor = 0

	default:
		return result, NewUnknownVersionIncrementError(string(incrementKind))
	}

	return result, nil
}

func Latest(tagList []string) (latest Version, found bool) {
	var (
		current Version
		isValid bool
	)

	for _, tag := range tagList {
		current, isValid = Parse(tag)
		if !isValid || (found && current.Compare(latest) <= 0) {
			continue
		}

		latest = current
		found = true
	}

	return latest, found
}

type UnknownVersionIncrementError struct {
	Kind string
}

func NewUnknownVersionIncrementError(kind string) (result *UnknownVersionIncrementError) {
	return &UnknownVersionIncrementError{Kind: kind}
}

func (e *UnknownVersionIncrementError) Error() (result string) {
	return fmt.Sprintf("unknown version increment %q: use patch, minor, or major", e.Kind)
}

type VersionOverflowError struct {
	Kind string
}

func NewVersionOverflowError(kind string) (result *VersionOverflowError) {
	return &VersionOverflowError{Kind: kind}
}

func (e *VersionOverflowError) Error() (result string) {
	return e.Kind + " version overflow"
}
