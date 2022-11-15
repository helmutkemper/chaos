package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
)

// SaVolumeList
//
// English:
//
//	List all docker volumes
//
//	 Output:
//	   list: docker volume list
//	   err: Standard error object
//
// Português:
//
//	Lista todos os volumes docker
//
//	 Saída:
//	   list: lista de volumes docker
//	   err: Objeto de erro padrão
func SaVolumeList() (list []types.Volume, err error) {
	var dockerSys = DockerSystem{}
	err = dockerSys.Init()
	if err != nil {
		util.TraceToLog()
		return
	}

	list, err = dockerSys.VolumeList()
	return
}
