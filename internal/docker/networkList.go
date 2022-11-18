package iotmakerdocker

import (
	"github.com/docker/docker/api/types"
)

func (el *DockerSystem) NetworkList() (
	netList []types.NetworkResource,
	err error,
) {

	netList, err = el.cli.NetworkList(el.ctx, types.NetworkListOptions{})
	return
}
