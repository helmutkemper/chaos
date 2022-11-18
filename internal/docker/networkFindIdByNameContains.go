package iotmakerdocker

import (
	"errors"
	"github.com/docker/docker/api/types"
	"strings"
)

func (el *DockerSystem) NetworkFindIdByNameContains(
	name string,
) (
	list []NameAndId,
	err error,
) {

	var listTmp []types.NetworkResource
	list = make([]NameAndId, 0)

	listTmp, err = el.NetworkList()
	if err != nil {
		return
	}

	for _, data := range listTmp {
		if strings.Contains(data.Name, name) == true {
			list = append(list, NameAndId{
				ID:   data.ID,
				Name: data.Name,
			})
		}
	}

	if len(list) == 0 {
		err = errors.New("network name not found")
	}

	return
}
