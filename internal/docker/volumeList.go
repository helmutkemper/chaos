package iotmakerdocker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	volumeTypes "github.com/docker/docker/api/types/volume"
)

func (el *DockerSystem) VolumeList() (
	volList []types.Volume,
	err error,
) {

	volList = make([]types.Volume, 0)

	var list volumeTypes.VolumeListOKBody
	list, err = el.cli.VolumeList(el.ctx, filters.Args{})
	if err != nil {
		return
	}

	for _, data := range list.Volumes {
		volList = append(volList, *data)
	}

	return
}
