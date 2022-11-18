package iotmakerdocker

import (
	"github.com/docker/docker/api/types"
)

func (el *DockerSystem) ImageRemove(
	id string,
	force,
	pruneChildren bool,
) (
	err error,
) {

	_, err = el.cli.ImageRemove(el.ctx, id, types.ImageRemoveOptions{
		Force:         force,
		PruneChildren: pruneChildren,
	})

	return err
}
