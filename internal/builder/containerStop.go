package builder

import (
	"time"
)

// ContainerStop (English): Stop a container by id
//
//	id: string container id
//
// ContainerStop (Português): Para um container por id
//
//	id: string container id
func (el *DockerSystem) ContainerStop(
	id string,
) (
	err error,
) {

	var timeout = time.Microsecond * 1000
	return el.cli.ContainerStop(el.ctx, id, &timeout)
}
