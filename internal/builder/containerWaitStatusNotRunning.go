package builder

// ContainerWaitStatusNotRunning (English): Waits until a container is in "not-running"
// status
//
//	id: string container id
//
// ContainerWaitStatusNotRunning (Português): Aguarda até o container entrar no estado de
// "not-running"
//
//	id: string container id
func (el *DockerSystem) ContainerWaitStatusNotRunning(
	id string,
) (
	err error,
) {

	wOk, wErr := el.cli.ContainerWait(el.ctx, id, "not-running")
	select {
	case <-wOk:
	case err = <-wErr:
	}
	return
}
