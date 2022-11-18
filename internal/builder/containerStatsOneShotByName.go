package builder

import (
	"errors"
	"github.com/docker/docker/api/types"
)

// ContainerStatisticsOneShotByName (English): Returns the performance information of the
// container in a timely manner by name
//
//	name: string container name
//
// ContainerStatisticsOneShotByName (Português): Retorna as informações de desempenho do
// container de forma pontual por nome
//
//	name: string container nome
func (el *DockerSystem) ContainerStatisticsOneShotByName(
	name string,
) (
	ret types.Stats,
	err error,
) {

	var list []types.Container
	var pass bool
	var id string

	list, err = el.ContainerListAll()
	if err != nil {
		return
	}

	for _, containerData := range list {
		for _, containerName := range containerData.Names {
			if containerName == name || containerName == "/"+name {
				pass = true
				id = containerData.ID
				break
			}
		}
	}

	if pass == false {
		return ret, errors.New("container name not found")
	}

	return el.ContainerStatisticsOneShot(id)
}
