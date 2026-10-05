package git

import (
	"context"
)

func (g Git) DeleteTag(ctx context.Context, tag string) (err error) {
	var HasTag bool

	HasTag, err = g.HasTag(ctx, tag)
	if err != nil || !HasTag {
		return err
	}

	_, err = g.RunCommand(ctx, "tag", "--delete", tag)
	if err != nil {
		return err
	}

	return nil
}
