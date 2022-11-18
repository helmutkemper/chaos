package iotmakerdocker

import (
	"github.com/docker/docker/api/types"
)

func (el *DockerSystem) NetworkInspect(
	id string,
) (
	inspect types.NetworkResource,
	err error,
) {

	inspect, err = el.cli.NetworkInspect(el.ctx, id, types.NetworkInspectOptions{})
	return inspect, err
}
