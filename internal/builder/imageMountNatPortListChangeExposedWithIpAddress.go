package builder

import (
	"github.com/docker/go-connections/nat"
)

// fixme: isto deveria ser privado?
func (el *DockerSystem) ImageMountNatPortListChangeExposedWithIpAddress(
	imageId,
	ipAddress string,
	currentPortList,
	changeToPortList []nat.Port,
) (
	nat.PortMap,
	error,
) {

	var err error
	var portList []nat.Port
	var ret nat.PortMap = make(map[nat.Port][]nat.PortBinding)

	portList, err = el.ImageListExposedPorts(imageId)
	if err != nil {
		return nat.PortMap{}, err
	}

	for _, port := range portList {
		inPort := ""
		for k, currPort := range currentPortList {
			if currPort.Port() == port.Port() && currPort.Proto() == port.Proto() {
				inPort = changeToPortList[k].Port()
				break
			}
		}

		ret[port] = []nat.PortBinding{
			{
				HostPort: inPort,
				HostIP:   ipAddress,
			},
		}
	}

	return ret, err
}
