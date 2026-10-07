package git_test

import (
	"context"
	"errors"
	"testing"

	"github.com/leprosus/tagit/internal/git"
)

func TestHistoryCancellation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	commits, err := git.New(t.TempDir()).GetHistory(ctx, "")
	if !errors.Is(err, context.Canceled) || len(commits) != 0 {
		t.Fatalf("commits=%v error=%v", commits, err)
	}
}
