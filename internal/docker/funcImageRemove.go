package docker

import (
	"github.com/helmutkemper/util"
)

// ImageRemove
//
// English:
//
//	Remove the image if there are no containers using the image
//
//	 Output:
//	   err: Standard error object
//
// Note:
//
//   - Remove all containers before use, including stopped containers
//
// Português:
//
//	Remove a imagem se não houver containers usando a imagem
//
//	 Saída:
//	   err: Objeto de erro padrão
//
// Nota:
//
//   - Remova todos os containers antes do uso, inclusive os containers parados
func (e *ContainerBuilder) ImageRemove() (err error) {
	err = e.dockerSys.ImageRemoveByName(e.imageName, false, false)
	if err != nil {
		util.TraceToLog()
	}
	return
}
