package docker

// ContainerFindIdByName
//
// Similar:
//
//	ContainerFindIdByName(), ContainerFindIdByNameContains()
//
// English:
//
//	Searches and returns the ID of the container, if it exists
//
//	 Input:
//	   name: Full name of the container.
//
//	 Output:
//	   id: container ID
//	   err: standard error object
//
// Português:
//
//	Procura e retorna o ID do container, caso o mesmo exista
//
//	 Entrada:
//	   name: Nome completo do container.
//
//	 Saída:
//	   id: ID do container
//	   err: Objeto de erro padrão
func (e *ContainerBuilder) ContainerFindIdByName(name string) (id string, err error) {
	e.dockerSys = DockerSystem{}
	err = e.dockerSys.Init()
	if err != nil {
		return
	}

	id, err = e.dockerSys.ContainerFindIdByName(name)
	if err != nil {
	}

	return
}
