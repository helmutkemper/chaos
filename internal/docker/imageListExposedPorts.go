package iotmakerdocker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"
)

// ImageListExposedPorts (English): List all ports / protocols published inside image
//
//	id: image ID
//
// Note: Similar functions: ImageListExposedNatPort(), ImageListExposedPortsByName(),
// ImageListExposedPorts()
//
// ImageListExposedPorts (Português): Lista todas as portas / protocolos publicadas
// dentro da imagem
//
//	id: ID da imagem
//
// Nota: funções similares: ImageListExposedNatPort(), ImageListExposedPortsByName(),
// ImageListExposedPorts()
func (el *DockerSystem) ImageListExposedPorts(
	id string,
) (
	portList []nat.Port,
	err error,
) {

	var imageData types.ImageInspect

	imageData, _, err = el.cli.ImageInspectWithRaw(el.ctx, id)
	if err != nil {
		return
	}
	for port := range imageData.ContainerConfig.ExposedPorts {
		portList = append(portList, port)
	}

	for port := range imageData.Config.ExposedPorts {
		portList = append(portList, port)
	}

	return
}
