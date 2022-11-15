package docker

// GetContainerIsStarted
//
// English:
//
//	Returns if the container was initialized after it was generated.
//
//	 Output:
//	   started: true for container initialized after generated
//
// Português:
//
//	Retorna se o container foi inicializado depois de gerado.
//
//	 Saída:
//	   started: true para container inicializado depois de gerado
func (e ContainerBuilder) GetContainerIsStarted() (started bool) {
	return e.startedAfterBuild
}
