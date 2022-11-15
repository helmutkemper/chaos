package docker

import (
	"github.com/docker/docker/api/types"
)

// ContainerStart (English): Start a container by id
//
//	id: string container id
//
// ContainerStart (Português): Inicia um container por id
//
//	id: string container id
func (el *DockerSystem) ContainerStart(
	id string,
) (
	err error,
) {

	return el.cli.ContainerStart(el.ctx, id, types.ContainerStartOptions{})
}
