package docker

// SetContainerName
//
// English:
//
//	Defines the name of the container
//
//	 Input:
//	   value: container name
//
// Português:
//
//	Define o nome do container
//
//	 Entrada:
//	   value: nome do container
func (e *ContainerBuilder) SetContainerName(value string) {
	e.containerName = value
}
