package docker

// GetInitialized
//
// English:
//
//	Returns if the container control object was initialized
//
//	 Output:
//	   initialized: true if the container control object was initialized
//
// Português:
//
//	Retorna se o objeto de controle do container foi inicializado
//
//	 Saída:
//	   inicializado: true caso o objeto de controle do container foi inicializado
func (e *ContainerBuilder) GetInitialized() (initialized bool) {
	return e.init
}
