package iotmakerdocker

import (
	"github.com/docker/docker/api/types"
)

func (el *DockerSystem) DockerInfo() (info types.Info, err error) {
	info, err = el.cli.Info(el.ctx)
	return
}
