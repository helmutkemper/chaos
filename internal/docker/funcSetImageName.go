package docker

import "strings"

// SetImageName
//
// English:
//
//	Defines the name of the image to be used or created
//
//	 Input:
//	   value: name of the image to be downloaded or created
//
// Note:
//
//   - the image name must have the version tag. E.g.: name:latest
//
// Português:
//
//	Define o nome da imagem a ser usada ou criada
//
//	 Entrada:
//	   value: noma da imagem a ser baixada ou criada
//
// Nota:
//
//   - o nome da imagem deve ter a tag de versão. Ex.: nome:latest
func (e *ContainerBuilder) SetImageName(value string) {
	if !strings.Contains(value, ":") {
		value += ":latest"
	}

	e.imageName = value
}
