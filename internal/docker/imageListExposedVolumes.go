package iotmakerdocker

import (
	"github.com/docker/docker/api/types"
)

// ImageListExposedVolumes (English): List all volumes exposed inside image file
// (dockerfile)
//
//	id: image ID
//
// ImageListExposedVolumes (Português): Lista todos os volumes expostos pela imagem
// (dockerfile)
//
//	id: ID da imagem
func (el *DockerSystem) ImageListExposedVolumes(
	id string,
) (
	list []string,
	err error,
) {

	var imageData types.ImageInspect
	list = make([]string, 0)

	imageData, _, err = el.cli.ImageInspectWithRaw(el.ctx, id)
	if err != nil {
		return []string{}, err
	}
	for volume := range imageData.ContainerConfig.Volumes {
		list = append(list, volume)
	}

	return
}
