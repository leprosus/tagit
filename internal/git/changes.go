package git

import "context"

// GetChanges returns the commits newest first without tag annotations.
func (g Git) GetChanges(ctx context.Context, tag string) (commits CommitInfoList, err error) {
	return g.getHistory(ctx, tag)
}
