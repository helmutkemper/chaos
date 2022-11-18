package builder

import (
	"time"
)

func (el *DockerSystem) ContainerRestartWithTimeout(
	id string,
	timeout time.Duration,
) (
	err error,
) {

	return el.cli.ContainerRestart(el.ctx, id, &timeout)
}
