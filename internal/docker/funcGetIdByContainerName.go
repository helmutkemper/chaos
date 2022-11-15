package docker

import (
	"github.com/helmutkemper/util"
)

// getIdByContainerName
//
// English:
//
//	Fills the container ID in the control object from the container name defined in SetContainerName()
//
//	 Output:
//	   err: Standard error object
//
// Português:
//
//	Preenche o ID do container no objeto de controle a partir do nome do container definido em
//	SetContainerName()
//
//	 Saída:
//	   err: Objeto de erro padrão
func (e *ContainerBuilder) getIdByContainerName() (err error) {
	e.containerID, err = e.dockerSys.ContainerFindIdByName(e.containerName)
	if err != nil {
		util.TraceToLog()
	}
	return
}
