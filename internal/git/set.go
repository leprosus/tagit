package git

import (
	"context"
)

func (g Git) CreateTag(ctx context.Context, tag string) (err error) {
	_, err = g.RunCommand(ctx, "tag", tag)
	if err != nil {
		return err
	}

	return nil
}
