package errors

import (
	"fmt"
	"strings"
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

type GitCommandError struct {
	ArgumentList []string
	Cause        error
	Output       string
}

func NewGitCommandError(argumentList []string, cause error, output string) (result *GitCommandError) {
	return &GitCommandError{ArgumentList: argumentList, Cause: cause, Output: output}
}

func (e *GitCommandError) Error() (result string) {
	return fmt.Sprintf(
		"git %s: %v: %s",
		strings.Join(e.ArgumentList, " "),
		e.Cause,
		e.Output,
	)
}

func (e *GitCommandError) Unwrap() (err error) {
	return e.Cause
}

type RemoteError struct {
	Remote string
	Cause  error
}

func NewRemoteError(remote string, cause error) (result *RemoteError) {
	return &RemoteError{Remote: remote, Cause: cause}
}

func (e *RemoteError) Error() (result string) {
	return fmt.Sprintf("remote %q: %v", e.Remote, e.Cause)
}

func (e *RemoteError) Unwrap() (err error) {
	return e.Cause
}
