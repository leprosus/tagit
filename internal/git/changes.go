package git

import "context"

// GetChanges returns the same commit details and range as GetHistory.
func (g Git) GetChanges(ctx context.Context, tag string) (commits CommitInfoList, err error) {
	return g.GetHistory(ctx, tag)
}
