package docker

import (
	"github.com/docker/docker/api/types"
)

// GetContainerInfo
//
// English:
//
//	Returns a series of information about the container.
//
//	 Output:
//	   info: Container information.
//	   err: Standard error object.
//
// Português:
//
//	Retorna uma séries de informações sobre o container.
//
//	 Saída:
//	   info: Informações sobre o container.
//	   err: Objeto padrão de erro.
func (e *ContainerBuilder) GetContainerInfo() (info types.Info, err error) {
	info, err = e.dockerSys.DockerInfo()
	return
}
