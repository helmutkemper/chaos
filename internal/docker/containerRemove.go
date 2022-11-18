package iotmakerdocker

import (
	"github.com/docker/docker/api/types"
)

// ContainerRemove (English): Remove a container by id
//
//	id: string container id
//	removeVolumes: bool remove container and volumes
//	removeLinks: bool remove container and links
//	force: bool force remove
//
// ContainerRemove (Português): Remove container por id
//
//	id: string container id
//	removeVolumes: bool remove o container e os volumes
//	removeLinks: bool remove o container e os links
//	force: bool força a emoção
func (el *DockerSystem) ContainerRemove(
	id string,
	removeVolumes,
	removeLinks,
	force bool,
) (
	err error,
) {

	return el.cli.ContainerRemove(
		el.ctx,
		id,
		types.ContainerRemoveOptions{
			RemoveVolumes: removeVolumes,
			RemoveLinks:   removeLinks,
			Force:         force,
		},
	)
}
