package git

import (
	"fmt"
	"io"
	"strings"
)

type CommitInfo struct {
	Commit  string
	Author  string
	Date    string
	Message string
}

type CommitInfoList []CommitInfo

func (c CommitInfo) Write(stdout io.Writer) (err error) {
	message, _, _ := strings.Cut(c.Message, "\n")

	message = strings.TrimSuffix(message, "\r")
	if message == "" {
		message = "(none)"
	}

	_, err = fmt.Fprintf(
		stdout,
		"Commit: %s\nAuthor: %s\nDate: %s\nMessage: %s\n",
		c.Commit,
		c.Author,
		c.Date,
		message,
	)

	return err
}

func (cl CommitInfoList) WriteCommitList(stdout io.Writer) (err error) {
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
