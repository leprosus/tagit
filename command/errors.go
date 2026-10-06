package command

import (
	"fmt"

	"github.com/leprosus/tagit/internal/version"
)

type UsageError struct{}

func NewUsageError() (result *UsageError) {
	return &UsageError{}
}

func (e *UsageError) Error() (result string) {
	result = "usage: tagit [patch|minor|major]"

	return result
}

type InvalidListLimitError struct {
	Value string
}

func NewInvalidListLimitError(value string) (result *InvalidListLimitError) {
	return &InvalidListLimitError{Value: value}
}

func (e *InvalidListLimitError) Error() (result string) {
	return fmt.Sprintf("invalid list limit %q: must be a positive integer", e.Value)
}

type NoVersionTagsError struct{}

func NewNoVersionTagsError() (result *NoVersionTagsError) {
	return &NoVersionTagsError{}
}

func (e *NoVersionTagsError) Error() (result string) {
	return "no local version tags found"
}

type UnknownVersionIncrementError = version.UnknownVersionIncrementError

func NewUnknownVersionIncrementError(kind string) (result *UnknownVersionIncrementError) {
	return version.NewUnknownVersionIncrementError(kind)
}

type VersionOverflowError = version.VersionOverflowError

func NewVersionOverflowError(kind string) (result *VersionOverflowError) {
	return version.NewVersionOverflowError(kind)
}
