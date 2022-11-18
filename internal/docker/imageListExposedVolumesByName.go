package iotmakerdocker

// list exposed volumes from image by name

// ImageListExposedVolumesByName (English): List all volumes exposed inside image file
// (dockerfile) by image name
//
//	name: name of image
//
// ImageListExposedVolumesByName (Português): Lista todos os volumes expostos pela imagem
// (dockerfile) pelo nome da imagem
//
//	name: nome da imagem
func (el *DockerSystem) ImageListExposedVolumesByName(
	name string,
) (
	list []string,
	err error,
) {

	var id string
	id, err = el.ImageFindIdByName(name)
	if err != nil {
		return nil, err
	}

	return el.ImageListExposedVolumes(id)
}
