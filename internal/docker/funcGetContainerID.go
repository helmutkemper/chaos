package docker

// GetContainerID
//
// English:
//
//	Returns the ID of the created container
//
//	 Output:
//	   ID: ID of the container
//
// Português:
//
//	Retorna o ID do container criado
//
//	 Saída:
//	   ID: ID do container
func (e *ContainerBuilder) GetContainerID() (ID string) {
	return e.containerID
}
