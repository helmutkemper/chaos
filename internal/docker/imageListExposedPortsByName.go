package docker

import "github.com/docker/go-connections/nat"

// ImageListExposedPortsByName (English): List all ports / protocols published inside image
// by image name
//
//	name: image name
//
// Note: Similar functions: ImageListExposedNatPort(), ImageListExposedPortsByName(),
// ImageListExposedPorts()
//
// ImageListExposedPortsByName (Português): Lista todas as portas / protocolos publicadas
// dentro da imagem pelo nome da imagem
//
//	name: nome da imagem
//
// Nota: funções similares: ImageListExposedNatPort(), ImageListExposedPortsByName(),
// ImageListExposedPorts()
func (el *DockerSystem) ImageListExposedPortsByName(
	name string,
) (
	portList []nat.Port,
	err error,
) {

	var id string
	id, err = el.ImageFindIdByName(name)
	if err != nil {
		return nil, err
	}

	portList, err = el.ImageListExposedPorts(id)

	return
}
