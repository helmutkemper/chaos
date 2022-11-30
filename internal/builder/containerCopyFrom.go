package builder

import (
	"github.com/docker/docker/api/types"
	"io"
)

// ContainerCopyFrom
//
// Copy from container to host
func (el *DockerSystem) ContainerCopyFrom(
	containerID string,
	sourcePath string,
) (
	reader io.ReadCloser,
	stats types.ContainerPathStat,
	err error,
) {

	reader, stats, err = el.cli.CopyFromContainer(el.ctx, containerID, sourcePath)
	return
}
