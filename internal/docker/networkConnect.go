package iotmakerdocker

import "github.com/docker/docker/api/types/network"

func (el *DockerSystem) NetworkConnect(
	networkID,
	containerID string,
	config *network.EndpointSettings,
) (
	err error,
) {

	return el.cli.NetworkConnect(el.ctx, networkID, containerID, config)
}
