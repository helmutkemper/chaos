package docker

import (
	"github.com/docker/docker/api/types"
)

// ContainerInspectByName (English): Inspect the container by name
//
//	name: string container name
//
// ContainerInspectByName (Português): Inspeciona o container pelo nome
//
//	name: string nome do container
func (el *DockerSystem) ContainerInspectByName(
	name string,
) (
	inspect types.ContainerJSON,
	err error,
) {

	var id string

	id, err = el.ContainerFindIdByName(name)
	if err != nil {
		return
	}

	inspect, err = el.ContainerInspect(id)

	return inspect, err
}
