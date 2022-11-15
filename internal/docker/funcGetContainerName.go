package docker

// GetContainerName
//
// English:
//
//	Returns container name
//
//	 Output:
//	   containerName: Container name
//
// Português:
//
//	Retorna o nome do container
//
//	 Saída:
//	   containerName: Nome do container
func (e *ContainerBuilder) GetContainerName() (containerName string) {
	return e.containerName
}
