package docker

import (
	"github.com/helmutkemper/util"
)

// ImageRemoveByName
//
// English:
//
//	Remove the image if there are no containers using the image
//
//	 Input:
//	   name: full image name
//
//	 Output:
//	   err: standard error object
//
// Note:
//
//   - Remove all containers before use, including stopped containers
//
// Português:
//
//	Remove a imagem se não houver containers usando a imagem
//
//	 Entrada:
//	   name: nome completo da imagem
//
//	 Saída:
//	   err: objeto de erro padrão
//
// Nota:
//
//   - Remova todos os containers antes do uso, inclusive os containers parados
func (e *ContainerBuilder) ImageRemoveByName(name string) (err error) {
	err = e.dockerSys.ImageRemoveByName(name, false, false)
	if err != nil {
		util.TraceToLog()
		return
	}

	return
}
