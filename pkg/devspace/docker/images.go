package docker

import (
	"context"
	dockertypes "github.com/docker/docker/api/types"
	"strings"

	"github.com/loft-sh/devspace/pkg/util/log"

	"github.com/docker/docker/api/types/filters"
)

// DeleteImageByName deletes an image by name
func (c *client) DeleteImageByName(ctx context.Context, imageName string, log log.Logger) ([]dockertypes.ImageDeleteResponseItem, error) {
	return c.DeleteImageByFilter(ctx, filters.NewArgs(filters.Arg("reference", strings.TrimSpace(imageName))), log)
}

// DeleteImageByFilter deletes an image by filter
func (c *client) DeleteImageByFilter(ctx context.Context, filter filters.Args, log log.Logger) ([]dockertypes.ImageDeleteResponseItem, error) {
	summaries, err := c.ImageList(ctx, dockertypes.ImageListOptions{
		Filters: filter,
	})
	if err != nil {
		return nil, err
	}

	responseItems := make([]dockertypes.ImageDeleteResponseItem, 0, 128)
	for _, summary := range summaries {
		deleteResponse, err := c.ImageRemove(ctx, summary.ID, dockertypes.ImageRemoveOptions{
			PruneChildren: true,
			Force:         true,
		})
		if err != nil {
			log.Warnf("%v", err)
		}

		if deleteResponse != nil {
			responseItems = append(responseItems, deleteResponse...)
		}
	}

	return responseItems, nil
}
