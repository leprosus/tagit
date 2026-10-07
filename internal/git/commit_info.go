package git

import (
	"fmt"
	"io"
	"strings"
)

type CommitInfo struct {
	Tag         string
	Commit      string
	ShortCommit string
	Author      string
	Date        string
	Message     string
}

type CommitInfoList []CommitInfo

func (c CommitInfo) Write(stdout io.Writer) (err error) {
	if c.Tag != "" {
		_, err = fmt.Fprintf(stdout, "Tag: %s\n", c.Tag)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(
		stdout,
		"Commit: %s\nAuthor: %s\nDate: %s\nMessage: %s\n",
		c.Commit,
		c.Author,
		c.Date,
		c.firstMessageLine(),
	)

	return err
}

func (cl CommitInfoList) Write(stdout io.Writer) (err error) {
	for i, commit := range cl {
		if i > 0 {
			_, err = io.WriteString(stdout, "\n")
			if err != nil {
				return err
			}
		}

		err = commit.Write(stdout)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c CommitInfo) WriteShort(stdout io.Writer) (err error) {
	_, err = fmt.Fprintf(
		stdout,
		"%s: %s\n",
		c.ShortCommit,
		c.firstMessageLine(),
	)

	return err
}

func (cl CommitInfoList) WriteShort(stdout io.Writer) (err error) {
	for _, commit := range cl {
		err = commit.WriteShort(stdout)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c CommitInfo) firstMessageLine() (message string) {
	message, _, _ = strings.Cut(c.Message, "\n")

	message = strings.TrimSuffix(message, "\r")
	if message == "" {
		return "(none)"
	}

	return message
}
