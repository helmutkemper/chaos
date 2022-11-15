package docker

// SetDockerfileBuilder
//
// English:
//
//	Defines a new object containing the builder of the dockerfile.
//
//	 Input:
//	   value: Object compatible with DockerfileAuto interface
//
// Note:
//
//   - Eee the DockerfileAuto interface for further instructions.
//
// Português:
//
//	Define um novo objeto contendo o construtor do arquivo dockerfile.
//
//	 Entrada:
//	   value: Objeto compatível com a interface DockerfileAuto
//
// Nota:
//
//   - Veja a interface DockerfileAuto para mais instruções.
func (e *ContainerBuilder) SetDockerfileBuilder(value DockerfileAuto) {
	e.autoDockerfile = value
}
