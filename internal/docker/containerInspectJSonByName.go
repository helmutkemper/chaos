package docker

// ContainerInspectJSonByName (English): Inspect the container by name e return a json
//
//	name: string container name
//
// ContainerInspectJSonByName (Português): Inspeciona o container pelo nome e retorna um
// json
//
//	name: string nome do container
func (el *DockerSystem) ContainerInspectJSonByName(
	name string,
) (
	inspect []byte,
	err error,
) {

	var id string

	id, err = el.ContainerFindIdByName(name)
	if err != nil {
		return
	}

	inspect, err = el.ContainerInspectJSon(id)

	return inspect, err
}
