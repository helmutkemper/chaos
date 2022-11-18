package iotmakerdocker

import (
	"errors"
	"github.com/docker/docker/api/types"
	"strings"
)

// ContainerFindIdByNameContains (English): Search by part of the container name and
// returns a list of NameAndId
//
//	name: string part of the container name
//
// ContainerFindIdByNameContains (Português): Procura por parte do nome do container e
// retorna uma lista de NameAndId
//
//	name: string parte do nome do container
func (el *DockerSystem) ContainerFindIdByNameContains(
	containsName string,
) (
	list []NameAndId,
	err error,
) {

	list = make([]NameAndId, 0)
	var listOfContainers []types.Container

	listOfContainers, err = el.ContainerListAll()
	if err != nil {
		return
	}

	for _, containerData := range listOfContainers {
		for _, containerName := range containerData.Names {
			if strings.Contains(containerName, containsName) == true {
				list = append(list, NameAndId{
					ID:   containerData.ID,
					Name: containerName,
				})
			}
		}
	}

	if len(list) == 0 {
		err = errors.New("container name not found")
	}

	return
}
