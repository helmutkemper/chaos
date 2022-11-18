package builder

// ContainerInspectJSon (English): Inspect the container by ID e return a json
//
//	name: string container ID
//
// ContainerInspectJSon (Português): Inspeciona o container pelo ID e retorna um json
//
//	name: string ID do container
func (el *DockerSystem) ContainerInspectJSon(
	id string,
) (
	inspect []byte,
	err error,
) {

	_, inspect, err = el.cli.ContainerInspectWithRaw(el.ctx, id, true)
	if err != nil {
		return
	}

	return
}
