package docker

// SetDockerfilePath
//
// English:
//
// Defines a Dockerfile to build the image.
//
// Português:
//
// Define um arquivo Dockerfile para construir a imagem.
func (e *ContainerBuilder) SetDockerfilePath(path string) (err error) {
	e.buildOptions.Dockerfile = path
	return
}
