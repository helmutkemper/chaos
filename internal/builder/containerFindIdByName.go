package builder

import (
	"errors"
	"github.com/docker/docker/api/types"
)

// ContainerFindIdByName (English): Searches for the container name and returns the ID of
// the container
//
//	name: string container name
//
// ContainerFindIdByName (Português): Procura pelo nome do container e retorna o ID do
// mesmo
//
//	name: string nome do container
func (el *DockerSystem) ContainerFindIdByName(
	name string,
) (
	ID string,
	err error,
) {

	var list []types.Container

	list, err = el.ContainerListAll()
	for _, containerData := range list {
		for _, containerName := range containerData.Names {
			if containerName == name || containerName == "/"+name {
				ID = containerData.ID
				return
			}
		}
	}

	err = errors.New("container name not found")
	return
}
