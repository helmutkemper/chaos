package docker

import (
	"github.com/helmutkemper/util"
)

// ImageFindIdByName
//
// English:
//
//	Find an image by name
//
//	 Input:
//	   name: image name
//
//	 Output:
//	   id: image ID
//	   err: default error object
//
// Português:
//
//	Encontra uma imagem pelo nome
//
//	 Input:
//	   name: nome da imagem
//
//	 Output:
//	   id: ID da imagem
//	   err: Objeto padrão de erro
func (e *ContainerBuilder) ImageFindIdByName(name string) (id string, err error) {
	e.dockerSys = DockerSystem{}
	err = e.dockerSys.Init()
	if err != nil {
		util.TraceToLog()
		return
	}

	id, err = e.dockerSys.ImageFindIdByName(name)
	if err != nil && err.Error() != "image name not found" {
		util.TraceToLog()
	}
	return
}
