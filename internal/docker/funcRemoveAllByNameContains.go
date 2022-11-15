package docker

import (
	"github.com/helmutkemper/util"
)

// RemoveAllByNameContains
//
// English:
//
//	Searches for networks, volumes, containers and images that contain the term defined in "value" in
//	the name, and tries to remove them from docker
//
//	 Input:
//	   value: part of the wanted name
//
//	 Output:
//	   err: Standard error object
//
// Português:
//
//	Procura por redes, volumes, container e imagens que contenham o termo definido em "value" no nome,
//	e tenta remover os mesmos do docker
//
//	 Entrada:
//	   value: parte do nome desejado
//
//	 Saída:
//	   err: Objeto de erro padrão
func (e *ContainerBuilder) RemoveAllByNameContains(value string) (err error) {
	e.containerID = ""
	err = e.dockerSys.RemoveAllByNameContains(value)
	if err != nil {
		util.TraceToLog()
	}
	return
}
