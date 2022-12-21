package builder

import (
	"context"
	"time"
)

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
	timeout time.Duration,
) (
	err error,
) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	wOk, wErr := el.cli.ContainerWait(ctx, id, "not-running")
	defer cancel()

	select {
	case <-wOk:
	case err = <-wErr:
	}
	return
}
