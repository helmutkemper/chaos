package builder

import (
	"time"
)

// ContainerStopAndRemove (English): Stop and removes a container by id
//
//	id: string container id
//
// ContainerStopAndRemove (Português): Para e remove um container por id
//
//	id: string container id
func (el *DockerSystem) ContainerStopAndRemove(
	id string,
	removeVolumes,
	removeLinks,
	force bool,
) (
	err error,
) {

	var timeout = time.Microsecond * 10000
	err = el.cli.ContainerStop(el.ctx, id, &timeout)
	if err != nil {
		return err
	}

	ok, notOk := el.cli.ContainerWait(el.ctx, id, "not-running")
	select {
	case <-ok:
		break
	case err = <-notOk:
		return err
	}

	time.Sleep(time.Second * 5)
	return el.ContainerRemove(id, removeVolumes, removeLinks, force)
}
