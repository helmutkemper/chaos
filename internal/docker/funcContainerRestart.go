package docker

import (
	"github.com/helmutkemper/util"
)

// ContainerRestart
//
// English:
//
// Restarts a container stopped by ContainerStop().
//
//	Output:
//	  err: standard error object
//
// Português:
//
//	Reinicia um container parado por ContainerStop().
//
//	 Saída:
//	   err: objeto de erro padrão
func (e *ContainerBuilder) ContainerRestart() (err error) {
	if e.containerID == "" {
		err = e.getIdByContainerName()
		if err != nil {
			util.TraceToLog()
			return
		}
	}

	err = e.dockerSys.ContainerRestart(e.containerID)
	if err != nil {
		util.TraceToLog()
	}
	return
}
