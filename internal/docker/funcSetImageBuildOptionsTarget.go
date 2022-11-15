package docker

// SetImageBuildOptionsTarget
//
// English:
//
//	Build the specified stage as defined inside the Dockerfile.
//
//	 Input:
//	   value: stage name
//
// Note:
//
//   - See the multi-stage build docs for details.
//     See https://docs.docker.com/develop/develop-images/multistage-build/
//
// Português:
//
//	Monta o container a partir do estágio definido no arquivo Dockerfile.
//
//	 Entrada:
//	   value: nome do estágio
//
// Nota:
//
//   - Veja a documentação de múltiplos estágios para mais detalhes.
//     See https://docs.docker.com/develop/develop-images/multistage-build/
func (e *ContainerBuilder) SetImageBuildOptionsTarget(value string) {
	e.buildOptions.Target = value
}
