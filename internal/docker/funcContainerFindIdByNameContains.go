package docker

// ContainerFindIdByNameContains
//
// Similar:
//
//	ContainerFindIdByName(), ContainerFindIdByNameContains()
//
// English:
//
//	Searches and returns the ID list of the container name
//
//	 Input:
//	   name: name of the container.
//
//	 Output:
//	   id: list of containers ID
//	   err: standard error object
//
// Português:
//
//	Procura e retorna uma lista de IDs de containers
//
//	 Entrada:
//	   name: Nome do container.
//
//	 Saída:
//	   id: lista de IDs dos containers
//	   err: Objeto de erro padrão
func (e *ContainerBuilder) ContainerFindIdByNameContains(containsName string) (list []NameAndId, err error) {
	list = make([]NameAndId, 0)

	e.dockerSys = DockerSystem{}
	err = e.dockerSys.Init()
	if err != nil {
		return
	}

	var receivedLis []NameAndId
	receivedLis, err = e.dockerSys.ContainerFindIdByNameContains(containsName)
	if err != nil {
		return
	}

	for _, elementInList := range receivedLis {
		list = append(list, NameAndId(elementInList))
	}

	return
}
