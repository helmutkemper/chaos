package iotmakerdocker

// ContainerWaitStatusRemoved (English): Waits until a container is in "removed"
// status
//
//	id: string container id
//
// ContainerWaitStatusRemoved (Português): Aguarda até o container entrar no estado de
// "removed"
//
//	id: string container id
func (el *DockerSystem) ContainerWaitStatusRemoved(
	id string,
) (
	err error,
) {

	wOk, wErr := el.cli.ContainerWait(el.ctx, id, "removed")
	select {
	case <-wOk:
	case err = <-wErr:
	}
	return
}
