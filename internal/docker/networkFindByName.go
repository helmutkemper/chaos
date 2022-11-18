package iotmakerdocker

import (
	"errors"
	"github.com/docker/docker/api/types"
)

func (el *DockerSystem) NetworkFindIdByName(
	name string,
) (
	id string,
	err error,
) {

	var list []types.NetworkResource

	list, err = el.NetworkList()
	if err != nil {
		return
	}

	for _, data := range list {
		if data.Name == name {
			id = data.ID
			return
		}
	}

	err = errors.New("network not found")

	return
}
