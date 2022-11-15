package docker

// SetContainerEntrypointToRunWhenStartingTheContainer
//
// English:
//
//	Entrypoint to run when stating the container
//
//	 Input:
//	   values: entrypoint. Eg.: docker run --entrypoint [new_command] [docker_image] [optional:value]
//
// Português:
//
//	Entrypoint a ser executado quando o container iniciar
//
//	 Entrada:
//	   values: entrypoint. Ex.: docker run --entrypoint [new_command] [docker_image] [optional:value]
func (e *ContainerBuilder) SetContainerEntrypointToRunWhenStartingTheContainer(values []string) {
	e.containerConfig.Entrypoint = values
}
