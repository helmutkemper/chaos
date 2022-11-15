package docker

import (
	"errors"
	"github.com/helmutkemper/util"
	"strings"
)

// verifyImageName
//
// English: check if the image name has the version tag
//
// Português: verifica se o nome da imagem tem a tag de versão
func (e *ContainerBuilder) verifyImageName() (err error) {
	if e.imageName == "" {
		util.TraceToLog()
		err = errors.New("image name is't set")
		return
	}

	if strings.Contains(e.imageName, ":") == false {
		util.TraceToLog()
		err = errors.New("image name must have a tag version. example: image_name:latest")
		return
	}

	return
}
