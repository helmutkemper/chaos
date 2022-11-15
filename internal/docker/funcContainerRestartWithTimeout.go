package docker

import (
	"github.com/helmutkemper/util"
	"time"
)

// ContainerRestartWithTimeout
//
// English:
//
//	Restarts a container stopped by ContainerStop().
//
//	 Input:
//	   timeout: timeout to restar container
//
//	 Output:
//	   err: standard error object
//
// Português:
//
//	Reinicia um container parado por ContainerStop().
//
//	 Entrada:
//	   timeout: tempo limite para reinício do container
//
//	 Saída:
//	   err: objeto de erro padrão
func (e *ContainerBuilder) ContainerRestartWithTimeout(timeout time.Duration) (err error) {
	if e.containerID == "" {
		err = e.getIdByContainerName()
		if err != nil {
			util.TraceToLog()
			return
		}
	}

	err = e.dockerSys.ContainerRestartWithTimeout(e.containerID, timeout)
	if err != nil {
		util.TraceToLog()
	}
	return
}
