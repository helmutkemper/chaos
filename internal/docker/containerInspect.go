package docker

import (
	"github.com/docker/docker/api/types"
)

// ContainerInspect (English): Inspect the container
//
//	id: string container ID
//
// ContainerInspect (Português): Inspeciona o container
//
//	id: string ID do container
func (el *DockerSystem) ContainerInspect(
	id string,
) (
	inspect types.ContainerJSON,
	err error,
) {

	inspect, err = el.cli.ContainerInspect(el.ctx, id)
	if err != nil {
		return
	}

	return
}
