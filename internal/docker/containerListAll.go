package docker

import "github.com/docker/docker/api/types"

// ContainerListAll (English): List all containers
//
// ContainerListAll (Português): Lista todos os containers
func (el *DockerSystem) ContainerListAll() (
	list []types.Container,
	err error,
) {

	list, err = el.cli.ContainerList(el.ctx, types.ContainerListOptions{All: true})
	return
}
