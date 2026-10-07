package git

import (
	"fmt"
	"strings"
)

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

type TagNotFoundError struct {
	Tag string
}

func NewTagNotFoundError(tag string) (result *TagNotFoundError) {
	return &TagNotFoundError{Tag: tag}
}

func (e *TagNotFoundError) Error() (result string) {
	return fmt.Sprintf("local tag %q not found", e.Tag)
}

type InvalidHistoryOutputError struct{}

func NewInvalidHistoryOutputError() (result *InvalidHistoryOutputError) {
	return &InvalidHistoryOutputError{}
}

func (*InvalidHistoryOutputError) Error() (result string) {
	return "invalid Git history output"
}
