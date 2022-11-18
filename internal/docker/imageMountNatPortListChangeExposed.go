package iotmakerdocker

import (
	"github.com/docker/go-connections/nat"
)

// fixme: isto deveria ser privado?
func (el *DockerSystem) ImageMountNatPortListChangeExposed(
	imageId string,
	currentPortList,
	changeToPortList []nat.Port,
) (
	nat.PortMap,
	error,
) {

	// fixme: ipAddress bug
	return el.ImageMountNatPortListChangeExposedWithIpAddress(imageId, "", currentPortList, changeToPortList)
}
