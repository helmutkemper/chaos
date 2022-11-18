package iotmakerdocker

import (
	"github.com/docker/docker/api/types"
	volumeTypes "github.com/docker/docker/api/types/volume"
)

func (el *DockerSystem) VolumeCreate(
	labels map[string]string,
	name string,
) (
	volume types.Volume,
	err error,
) {

	volume, err = el.cli.VolumeCreate(el.ctx, volumeTypes.VolumeCreateBody{Labels: labels, Name: name})

	return
}
