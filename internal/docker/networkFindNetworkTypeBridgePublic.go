package iotmakerdocker

import (
	"errors"
	"github.com/docker/docker/api/types"
)

func (el *DockerSystem) NetworkFindNetworkTypeBridgePublic() (
	inspect types.NetworkResource,
	err error,
) {

	var list []types.NetworkResource
	var netDriveToFind = KNetworkDriveBridge

	list, err = el.NetworkList()
	if err != nil {
		return
	}

	for _, net := range list {
		if net.Driver == netDriveToFind.String() {
			inspect, err = el.NetworkInspect(net.ID)
			if err != nil {
				return
			}

			for k, v := range inspect.Options {
				if k == "com.docker.network.bridge.default_bridge" && v == "true" {
					return
				}
			}
		}
	}

	err = errors.New("network type bridge not found")
	inspect = types.NetworkResource{}

	return
}
