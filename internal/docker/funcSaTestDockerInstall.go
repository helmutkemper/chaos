package docker

import (
	"github.com/helmutkemper/util"
)

// SaTestDockerInstall
//
// English:
//
//	Test if docker is responding correctly
//
//	 Output:
//	   err: Standard error object
//
// Português:
//
//	Testa se o docker está respondendo de forma correta
//
//	 Saída:
//	   err: Standard error object
func SaTestDockerInstall() (err error) {
	var dockerSys = DockerSystem{}
	err = dockerSys.Init()
	if err != nil {
		util.TraceToLog()
		return
	}

	_, err = dockerSys.ImageList()
	return
}
