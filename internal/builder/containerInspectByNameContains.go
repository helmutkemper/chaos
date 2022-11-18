package builder

import (
	"github.com/docker/docker/api/types"
)

// ContainerInspectByNameContains (English): Inspects all containers whose part of the name
// contains the search term
//
//	name: string search term
//
// ContainerInspectByNameContains (Português): Inspeciona todos os containers com cuja parte do
// nome contém o termo procurado
//
//	name: string termo procurado
func (el *DockerSystem) ContainerInspectByNameContains(
	searchTerm string,
) (
	list []types.ContainerJSON,
	err error,
) {

	list = make([]types.ContainerJSON, 0)
	var inspect types.ContainerJSON
	var listOfContainers []NameAndId

	listOfContainers, err = el.ContainerFindIdByNameContains(searchTerm)
	if err != nil {
		return
	}

	for _, v := range listOfContainers {
		inspect, err = el.ContainerInspect(v.ID)
		if err != nil {
			return
		}

		list = append(list, inspect)
	}

	return
}
