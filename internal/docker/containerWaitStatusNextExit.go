package docker

// ContainerWaitStatusNextExit (English): Waits until a container is in "next-exit"
// status
//
//	id: string container id
//
// ContainerWaitStatusNextExit (Português): Aguarda até o container entrar no estado de
// "next-exit"
//
//	id: string container id
func (el *DockerSystem) ContainerWaitStatusNextExit(
	id string,
) (
	err error,
) {

	wOk, wErr := el.cli.ContainerWait(el.ctx, id, "next-exit")
	select {
	case <-wOk:
	case err = <-wErr:
	}
	return
}
