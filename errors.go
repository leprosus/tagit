package main

import (
	"fmt"
	"strings"
)

type usageError struct{}

func newUsageError() (result *usageError) {
	return &usageError{}
}

func (e *usageError) Error() (result string) {
	result = "usage: tagit [patch|minor|major]"

	return result
}

type unknownVersionIncrementError struct {
	kind kind
}

func newUnknownVersionIncrementError(kind kind) (result *unknownVersionIncrementError) {
	return &unknownVersionIncrementError{kind: kind}
}

func (e *unknownVersionIncrementError) Error() (result string) {
	return fmt.Sprintf("unknown version increment %q: use patch, minor, or major", e.kind)
}

type versionOverflowError struct {
	kind kind
}

func newVersionOverflowError(kind kind) (result *versionOverflowError) {
	return &versionOverflowError{kind: kind}
}

func (e *versionOverflowError) Error() (result string) {
	return fmt.Sprintf("%s version overflow", e.kind)
}

type gitCommandError struct {
	argumentList []string
	cause        error
	output       string
}

func newGitCommandError(argumentList []string, cause error, output string) (result *gitCommandError) {
	return &gitCommandError{argumentList: argumentList, cause: cause, output: output}
}

func (e *gitCommandError) Error() (result string) {
	return fmt.Sprintf("git %s: %v: %s", strings.Join(e.argumentList, " "), e.cause, e.output)
}

func (e *gitCommandError) Unwrap() (err error) {
	return e.cause
}
