package git_test

import (
	"context"
	"errors"
	"testing"

	"github.com/leprosus/tagit/internal/git"
)

func TestGetLatestVersionPreservesGitFailure(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	latest, found, err := git.New(t.TempDir()).GetLatestVersion(ctx)
	if found || latest.String() != "v0.0.0" || !errors.Is(err, context.Canceled) {
		t.Fatalf(
			"latest=%v found=%t error=%v, want original cancellation error",
			latest,
			found,
			err,
		)
	}
}
