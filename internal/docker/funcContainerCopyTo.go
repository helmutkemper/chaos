package docker

import (
	"github.com/helmutkemper/util"
	"io"
	"os"
)

// ContainerCopyTo
//
// Português:
//
//	Copia um arquivo contido no computador local para dentro do container
//
//	 Entrada:
//	   hostPathList: lista de arquivos a serem salvos no computador hospedeiro (caminho + nome do
//	     arquivo)
//	   containerPathList: lista de arquivos contidos no container (apenas o caminho)
//
//	 Saída:
//	   err: Objeto de erro padrão
//
// English:
//
//	Copy a file contained on the local computer into the container
//
//	 Input:
//	   hostPathList: list of files to be saved on the host computer (path + filename)
//	   containerPathList: list of files contained in the container (path only)
//
//	 Output:
//	   err: standard error object
func (e *ContainerBuilder) ContainerCopyTo(
	hostPathList []string,
	containerPathList []string,
) (
	err error,
) {
	if e.containerID == "" {
		err = e.getIdByContainerName()
		if err != nil {
			return
		}
	}

	var content io.Reader
	for k, destinationPath := range hostPathList {
		content, err = os.Open(destinationPath)
		if err != nil {
			util.TraceToLog()
			return
		}

		err = e.dockerSys.ContainerCopyTo(
			e.containerID,
			containerPathList[k],
			content,
		)
		if err != nil {
			util.TraceToLog()
			return
		}
	}

	return
}
