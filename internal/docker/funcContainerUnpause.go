package docker

import (
	"github.com/helmutkemper/util"
)

// ContainerUnpause
//
// English:
//
//	Remove the pause from the previously paused container with the container.Pause() command
//
//	 Output:
//	   err: Standard error object.
//
// Português:
//
//	Remove a pausa do container previamente pausado com o comando container.Pause()
//
//	 Saída:
//	   err: Objeto de erro padrão.
func (e *ContainerBuilder) ContainerUnpause() (err error) {
	if e.containerID == "" {
		err = e.getIdByContainerName()
		if err != nil {
			util.TraceToLog()
			return
		}
	}

	err = e.dockerSys.ContainerUnpause(e.containerID)
	if err != nil {
		util.TraceToLog()
	}
	return
}
