package builder

import (
	"github.com/docker/docker/api/types"
	"io"
)

// ContainerCopyTo
//
// Copy to host from container
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
