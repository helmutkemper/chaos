package builder

import (
	"github.com/docker/go-connections/nat"
)

// Mount nat por list by image config

// ImageListExposedNatPort (English): List all ports / protocols published inside image
//
//	name: image name
//
// Note: Similar functions: ImageListExposedNatPort(), ImageListExposedPortsByName(),
// ImageListExposedPorts()
//
// ImageListExposedNatPort (Português): Lista todas as portas / protocolos publicadas
// dentro da imagem
//
//	name: nome da imagem
//
// Nota: funções similares: ImageListExposedNatPort(), ImageListExposedPortsByName(),
// ImageListExposedPorts()
func (el *DockerSystem) ImageListExposedNatPort(
	imageId string,
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
		ret[port] = []nat.PortBinding{
			{
				HostPort: port.Port() + "/" + port.Proto(),
			},
		}
	}

	return ret, err
}
