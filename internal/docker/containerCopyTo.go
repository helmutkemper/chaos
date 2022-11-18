package iotmakerdocker

import (
	"github.com/docker/docker/api/types"
	"io"
)

func (el *DockerSystem) ContainerCopyTo(
	containerID string,
	destinationPath string,
	content io.Reader,
) (
	err error,
) {

	err = el.cli.CopyToContainer(el.ctx, containerID, destinationPath, content, types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
		CopyUIDGID:                false,
	})

	return
}
