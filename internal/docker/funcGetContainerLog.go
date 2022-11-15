package docker

import (
	"github.com/helmutkemper/util"
)

// GetContainerLog
//
// English:
//
//	Returns the current standard output of the container.
//
//	 Output:
//	   log: Texts contained in the container's standard output
//	   err: Standard error object
//
// Português:
//
//	Retorna a saída padrão atual do container.
//
//	 Saída:
//	   log: Textos contidos na saída padrão do container
//	   err: Objeto de erro padrão
func (e *ContainerBuilder) GetContainerLog() (log []byte, err error) {
	if e.containerID == "" {
		err = e.getIdByContainerName()
		if err != nil {
			util.TraceToLog()
			return
		}
	}

	log, err = e.dockerSys.ContainerLogs(e.containerID)
	if err != nil {
		util.TraceToLog()
	}
	return
}
