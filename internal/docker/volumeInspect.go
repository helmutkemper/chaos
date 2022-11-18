package iotmakerdocker

import (
	"github.com/docker/docker/api/types"
)

func (el *DockerSystem) VolumeInspect(
	ID string,
) (
	inspect types.Volume,
	err error,
) {

	inspect, err = el.cli.VolumeInspect(el.ctx, ID)
	return
}
